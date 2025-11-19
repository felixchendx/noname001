package pubsub

import (
	"context"
	
	"noname001/logging"

	dilemmaComm "noname001/dilemma/comm"
)

type DevicePublisherParams struct{
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
	
	CollectorServerHost string
}
type DevicePublisher struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	pubClient *dilemmaComm.PublisherClient
}

func NewDevicePublisher(params *DevicePublisherParams) (*DevicePublisher, error) {
	var err error

	pub := &DevicePublisher{}
	pub.context, pub.cancel = context.WithCancel(params.Context)
	pub.logger = params.Logger
	pub.logPrefix = params.LogPrefix + ".dvc.pub"

	pub.pubClient, err = dilemmaComm.NewPublisherClient(dilemmaComm.PublisherClientConfig{
		Context: pub.context,
		CollectorServerHost: params.CollectorServerHost,
	})
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func (pub *DevicePublisher) Connect() (error) {
	var err error

	err = pub.pubClient.Connect()

	if err != nil {
		return err
	}

	return nil
}

func (pub *DevicePublisher) Disconnect() (error) {
	pub.cancel()
	// pub.pubClient.Disconnect()
	return nil
}
