package comm

import (
	"context"

	"noname001/logging"

	"noname001/node/commconf"

	deviceMdp    "noname001/app/node/comm/device/mdpservice"
	devicePubsub "noname001/app/node/comm/device/pubsub"
)

type CommBundleParams struct {
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
}

type CommBundle struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	DeviceProvider  *deviceMdp.DeviceProvider
	DevicePublisher *devicePubsub.DevicePublisher
}

func NewCommBundle(params *CommBundleParams) (*CommBundle, error) {
	var err error
	
	commBundle := &CommBundle{}
	commBundle.context, commBundle.cancel = context.WithCancel(params.Context)
	commBundle.logger = params.Logger
	commBundle.logPrefix = params.LogPrefix + ".comm"

	commBundle.DeviceProvider, err = deviceMdp.NewDeviceProvider(&deviceMdp.DeviceProviderParams{
		Context: commBundle.context,
		Logger: commBundle.logger,
		LogPrefix: commBundle.logPrefix,

		TempNodeID: commconf.ID(),
		BrokerHost: commconf.BrokerServerHost(),
		DataHandler: nil,

		Verbose: commconf.CommVerbose(),
		RetryBackoff: deviceMdp.DEFAULT_RETRY_BACKOFF,
	})
	if err != nil {
		return nil, err
	}

	commBundle.DevicePublisher, err = devicePubsub.NewDevicePublisher(&devicePubsub.DevicePublisherParams{
		Context: commBundle.context,
		Logger : commBundle.logger, LogPrefix: commBundle.logPrefix,

		CollectorServerHost: commconf.CollectorServerHost(),
	})
	if err != nil {
		return nil, err
	}

	return commBundle, nil
}

func (commBundle *CommBundle) Connect() (err error) {
	err = commBundle.DeviceProvider.Connect()
	if err != nil { return }

	err = commBundle.DevicePublisher.Connect()
	if err != nil { return }

	return
}

func (commBundle *CommBundle) Disconnect() (err error) {
	err = commBundle.DevicePublisher.Disconnect()
	if err != nil { return }

	err = commBundle.DeviceProvider.Disconnect()
	if err != nil { return }

	return
}
