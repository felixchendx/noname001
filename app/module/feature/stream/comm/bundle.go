package comm

import (
	"context"

	"noname001/logging"

	"noname001/node/commconf"

	streamMdp    "noname001/app/node/comm/stream/mdpservice"
	streamPubsub "noname001/app/node/comm/stream/pubsub"
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

	StreamProvider  *streamMdp.StreamProvider
	StreamPublisher *streamPubsub.StreamPublisher
}

func NewCommBundle(params *CommBundleParams) (*CommBundle, error) {
	var err error
	
	commBundle := &CommBundle{}
	commBundle.context, commBundle.cancel = context.WithCancel(params.Context)
	commBundle.logger = params.Logger
	commBundle.logPrefix = params.LogPrefix + ".comm"

	commBundle.StreamProvider, err = streamMdp.NewStreamProvider(&streamMdp.StreamProviderParams{
		Context: commBundle.context,
		Logger: commBundle.logger,
		LogPrefix: commBundle.logPrefix,

		TempNodeID: commconf.ID(),
		BrokerHost: commconf.BrokerServerHost(),
		DataHandler: nil,

		Verbose: commconf.CommVerbose(),
		RetryBackoff: streamMdp.DEFAULT_RETRY_BACKOFF,
	})
	if err != nil {
		return nil, err
	}

	commBundle.StreamPublisher, err = streamPubsub.NewStreamPublisher(&streamPubsub.StreamPublisherParams{
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
	err = commBundle.StreamProvider.Connect()
	if err != nil { return }

	err = commBundle.StreamPublisher.Connect()
	if err != nil { return }

	return
}

func (commBundle *CommBundle) Disconnect() (err error) {
	err = commBundle.StreamPublisher.Disconnect()
	if err != nil { return }

	err = commBundle.StreamProvider.Disconnect()
	if err != nil { return }

	return
}
