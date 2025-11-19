package ws

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/fasthttp/websocket"

	"noname001/logging"
)

type WSClient struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	id     string
	wsConn *websocket.Conn

	// TODO: metadata mainly for admin area
	//       use map
	// webSessionID string
	// path + queryargs

	sendChan chan []byte
	recvChan chan []byte

	isDisconnecting bool // TODO: infer from context, how ? read doc
}

func (client *WSClient) Send(msg []byte) (error) {
	if client.isDisconnecting {
		err := fmt.Errorf("aborted send due to disconnecting...")
		client.logger.Warnf("%s: wscli-%s, ", client.logPrefix, client.id, err.Error())
		return err
	}

	client.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
	err := client.wsConn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		client.logger.Errorf("%s: wscli-%s, send err: %s", client.logPrefix, client.id, err.Error())
		return err
	}

	return nil
}

func (client *WSClient) ReceiverChannel() (chan []byte) {
	return client.recvChan
}

func (client *WSClient) disconnect() {
	client.isDisconnecting = true
	client.cancel()
}

func (client *WSClient) sender() {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()

		client.cleanup()
	}()

	senderLoop:
	for {
		select {
		case <- client.context.Done():
			err := client.wsConn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			_ = err
			break senderLoop

		case <- pingTicker.C:
			// for now, fail to ping immediately kill the connection
			// let the client re-establish new ws conn when the client can reach server
			client.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			err := client.wsConn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				client.logger.Warnf("%s: wscli-%s, ping err: %s", client.logPrefix, client.id, err.Error())
				break senderLoop
			}
		}
	}
}

func (client *WSClient) receiver() {
	defer func() {
		client.cleanup()
	}()

	client.wsConn.SetReadLimit(maxMessageSize)
	client.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	client.wsConn.SetPongHandler(func(string) error {
		client.wsConn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// exiting this loop will close the conn
	receiverLoop:
	for {
		if client.wsConn == nil {
			break receiverLoop
		}

		msgType, msg, err := client.wsConn.ReadMessage()
		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	client.logger.Debugf("TEMP:receiverLoop:XXX %s", err.Error())
			// 	break receiverLoop
			// }
				
			// TODO: err filtering
			// TODO: ws stat 1000 (normal)
			// TODO: ws stat 1001 (going away)

			client.logger.Errorf("%s: wscli-%s, recv err: %s", client.logPrefix, client.id, err.Error())
			break receiverLoop
		}

		switch msgType {
		case websocket.TextMessage:
			msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
			client.recvChan <- msg

		// TODO: read doc, more msgType handling

		default:
			client.logger.Warnf("%s: wscli-%s, unhandled msgType %s, %s", msgType, msg)
		}
	}
}

func (client *WSClient) cleanup() {
	// use wsConn as makeshift cleanup indicator
	// to avoid closing closed channel
	if client.wsConn != nil {
		client.wsConn.Close()
		client.wsConn = nil

		close(client.sendChan)
		close(client.recvChan)
	}
}
