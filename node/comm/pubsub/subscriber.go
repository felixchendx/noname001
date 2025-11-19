package pubsub

import (
	"context"
	"strings"

	"noname001/logging"

	dilemmaComm "noname001/dilemma/comm"
)

type NodeSubscriberParams struct {
	Context             context.Context
	Logger              *logging.WrappedLogger
	LogPrefix           string
	
	SnapshotServerHost  string
	PublisherServerHost string
}
type NodeSubscriber struct {
	context         context.Context
	cancel          context.CancelFunc
	logger          *logging.WrappedLogger
	logPrefix       string

	rawMessageChan  chan string
	subClient       *dilemmaComm.SubscriberClient

	// TODO: typing and assertion shenanigans
	nodeLivenessSubscriptions map[string]*NodeLivenessSubscription
	nodeEventSubscriptions map[string]*NodeEventSubscription
}

func NewNodeSubscriber(params *NodeSubscriberParams) (*NodeSubscriber, error) {
	var err error
	
	sub := &NodeSubscriber{}
	sub.context, sub.cancel = context.WithCancel(params.Context)
	sub.logger = params.Logger
	sub.logPrefix = params.LogPrefix + ".node.sub"

	sub.rawMessageChan = make(chan string)
	sub.subClient, err = dilemmaComm.NewSubscriber(dilemmaComm.SubscriberClientConfig{
		Context: sub.context,
		SnapshotServerHost: params.SnapshotServerHost,
		PublisherServerHost: params.PublisherServerHost,
		Subtree: SUBTREE__NODE_BASE,

		DataChannel: sub.rawMessageChan,
	})
	if err != nil {
		return nil, err
	}

	sub.nodeLivenessSubscriptions = make(map[string]*NodeLivenessSubscription)
	sub.nodeEventSubscriptions = make(map[string]*NodeEventSubscription)

	return sub, nil
}

func (sub *NodeSubscriber) Connect() (error) {
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
				case HEADER__EVENT   : sub.receiveNodeEvent(msgParts[1])
				case HEADER__LIVENESS: sub.receiveNodeLiveness(msgParts[1])
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

func (sub *NodeSubscriber) Disconnect() (error) {
	// TODO: terminate subscriptions

	sub.cancel()
	// sub.subClient.Disconnect()

	return nil
}
