package pubsub

import (
	"context"
	
	"noname001/logging"

	dilemmaComm "noname001/dilemma/comm"
)

type StreamPublisherParams struct{
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
	
	CollectorServerHost string
}
type StreamPublisher struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	pubClient *dilemmaComm.PublisherClient
}

func NewStreamPublisher(params *StreamPublisherParams) (*StreamPublisher, error) {
	var err error

	pub := &StreamPublisher{}
	pub.context, pub.cancel = context.WithCancel(params.Context)
	pub.logger = params.Logger
	pub.logPrefix = params.LogPrefix + ".stream.pub"

	pub.pubClient, err = dilemmaComm.NewPublisherClient(dilemmaComm.PublisherClientConfig{
		Context: pub.context,
		CollectorServerHost: params.CollectorServerHost,
	})
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func (pub *StreamPublisher) Connect() (error) {
	var err error

	err = pub.pubClient.Connect()

	if err != nil {
		return err
	}

	return nil
}

func (pub *StreamPublisher) Disconnect() (error) {
	pub.cancel()
	// pub.pubClient.Disconnect()
	return nil
}
