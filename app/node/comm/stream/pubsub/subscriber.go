package pubsub

import (
	"context"
	"strings"

	"noname001/logging"

	dilemmaComm "noname001/dilemma/comm"
)

type StreamSubscriberParams struct {
	Context             context.Context
	Logger              *logging.WrappedLogger
	LogPrefix           string
	
	SnapshotServerHost  string
	PublisherServerHost string
}
type StreamSubscriber struct {
	context        context.Context
	cancel         context.CancelFunc
	logger         *logging.WrappedLogger
	logPrefix      string

	rawMessageChan chan string
	subClient      *dilemmaComm.SubscriberClient

	streamEventSubscriptions map[string]*StreamEventSubscription
}

func NewStreamSubscriber(params *StreamSubscriberParams) (*StreamSubscriber, error) {
	var err error
	
	sub := &StreamSubscriber{}
	sub.context, sub.cancel = context.WithCancel(params.Context)
	sub.logger = params.Logger
	sub.logPrefix = params.LogPrefix + ".stream.sub"

	sub.rawMessageChan = make(chan string)
	sub.subClient, err = dilemmaComm.NewSubscriber(dilemmaComm.SubscriberClientConfig{
		Context: sub.context,
		SnapshotServerHost: params.SnapshotServerHost,
		PublisherServerHost: params.PublisherServerHost,
		Subtree: SUBTREE__STREAM_BASE,

		DataChannel: sub.rawMessageChan,
	})
	if err != nil {
		return nil, err
	}

	sub.streamEventSubscriptions = make(map[string]*StreamEventSubscription)

	return sub, nil
}

func (sub *StreamSubscriber) Connect() (error) {
	var err error

	err = sub.subClient.Connect()

	if err != nil {
		return err
	}

	// TODO: benchmark throughput of this setup (single goroutine)
	go func() {
		for {
			select {
			case <- sub.context.Done():
				break

			case rawMessage := <- sub.rawMessageChan:
				if rawMessage == "" {
					// reconfirm the purpose of publishing delete
					continue
				}

				msgParts := strings.Split(rawMessage, HEADER_PAYLOAD_DELIM)
				if len(msgParts) != 2 {
					// TODO:
					sub.logger.Warnf("%s: broken msg ? %s", sub.logPrefix, rawMessage)
					continue
				}

				switch msgParts[0] {
				case HEADER__EVENT: sub.receiveStreamEvent(msgParts[1])
				default:
					// TODO
					sub.logger.Warnf("%s: unknown header %s", sub.logPrefix, msgParts[0])
					continue
				}
			}
		}
	}()

	return nil
}

func (sub *StreamSubscriber) Disconnect() (error) {
	// TODO: terminate subscriptions

	sub.cancel()
	// sub.subClient.Disconnect()

	return nil
}
