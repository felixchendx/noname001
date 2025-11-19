package ws

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/fasthttp/websocket"

	"noname001/logging"
)

// same as lapi stuffs
// these ws stuffs are intended only for simple internal usage
// where the programmer has access to codes of both provider and consumer
// is not suitable for setup where be and fe are separate entities

const (
	// TODO: check configs
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte("\n")
	space   = []byte(" ")

	wsUpgrader = websocket.FastHTTPUpgrader{
		ReadBufferSize   : 1024,
		WriteBufferSize  : 1024,
		EnableCompression: false,
		// CheckOrigin      : func(r *http.Request) bool {
		// 	// TODO
		// 	return true
		// },
	}
)

type WSHubParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger
	LogPrefix     string
}
type WSHub struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	clients map[string]*WSClient
}

func NewWSHub(params *WSHubParams) (*WSHub) {
	hub := &WSHub{}
	hub.context, hub.cancel = context.WithCancel(params.ParentContext)
	hub.logger = params.Logger
	hub.logPrefix = params.LogPrefix + ".wshub"
	hub.clients = make(map[string]*WSClient)

	return hub
}

func (hub *WSHub) UpgradeToWebsocketConnection(ctx *fasthttp.RequestCtx) (*WSClient, error) {
	client := &WSClient{}
	client.context, client.cancel = context.WithCancel(hub.context)
	client.logger = hub.logger
	client.logPrefix = hub.logPrefix + ".wscli"
	client.id = uuid.New().String()
	client.wsConn = nil
	client.sendChan = make(chan []byte)
	client.recvChan = make(chan []byte)

	hub.clients[client.id] = client


	err := wsUpgrader.Upgrade(ctx, func(wsConn *websocket.Conn) {
		client.wsConn = wsConn

		go client.sender()
		client.receiver()
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (hub *WSHub) CloseWebsocketConnection(client *WSClient) {
	client.disconnect()
}
