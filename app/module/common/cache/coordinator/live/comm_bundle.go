package live

import (
	"noname001/node/commconf"

	tmphubMdp "noname001/hub/comm/tmphub/mdpservice"

	nodeMdp    "noname001/node/comm/mdpservice"
	nodePubsub "noname001/node/comm/pubsub"

	deviceMdp    "noname001/app/node/comm/device/mdpservice"
	devicePubsub "noname001/app/node/comm/device/pubsub"

	streamMdp    "noname001/app/node/comm/stream/mdpservice"
	streamPubsub "noname001/app/node/comm/stream/pubsub"
)

type t_commBundle struct {
	tmpHubReceiver *tmphubMdp.TmpHubReceiver

	nodeReceiver   *nodeMdp.NodeReceiver
	nodeSubscriber *nodePubsub.NodeSubscriber

	deviceReceiver   *deviceMdp.DeviceReceiver
	deviceSubscriber *devicePubsub.DeviceSubscriber

	streamReceiver   *streamMdp.StreamReceiver
	streamSubscriber *streamPubsub.StreamSubscriber
}

func (lc *LiveCache) newCommBundle() (*t_commBundle, error) {
	var (
		_commBundle       = &t_commBundle{}
		err         error
	)

	_commBundle.tmpHubReceiver, err = tmphubMdp.NewTmpHubReceiver(&tmphubMdp.TmpHubReceiverParams{
		Logger: lc.logger, LogPrefix: lc.logPrefix,

		BrokerHost: commconf.BrokerServerHost(),

		Verbose: commconf.CommVerbose(),
	})
	if err != nil {
		return nil, err
	}

	_commBundle.nodeReceiver, err = nodeMdp.NewNodeReceiver(&nodeMdp.NodeReceiverParams{
		Logger: lc.logger, LogPrefix: lc.logPrefix,

		BrokerHost: commconf.BrokerServerHost(),

		Verbose: commconf.CommVerbose(),
	})
	if err != nil {
		return nil, err
	}

	_commBundle.nodeSubscriber, err = nodePubsub.NewNodeSubscriber(&nodePubsub.NodeSubscriberParams{
		Context: lc.context,
		Logger: lc.logger, LogPrefix: lc.logPrefix,

		SnapshotServerHost: commconf.SnapshotServerHost(),
		PublisherServerHost: commconf.PublisherServerHost(),
	})
	if err != nil {
		return nil, err
	}

	_commBundle.deviceReceiver, err = deviceMdp.NewDeviceReceiver(&deviceMdp.DeviceReceiverParams{
		Logger: lc.logger, LogPrefix: lc.logPrefix,

		BrokerHost: commconf.BrokerServerHost(),

		Verbose: commconf.CommVerbose(),
	})
	if err != nil {
		return nil, err
	}

	_commBundle.deviceSubscriber, err = devicePubsub.NewDeviceSubscriber(&devicePubsub.DeviceSubscriberParams{
		Context: lc.context,
		Logger: lc.logger, LogPrefix: lc.logPrefix,

		SnapshotServerHost: commconf.SnapshotServerHost(),
		PublisherServerHost: commconf.PublisherServerHost(),
	})
	if err != nil {
		return nil, err
	}

	_commBundle.streamReceiver, err = streamMdp.NewStreamReceiver(&streamMdp.StreamReceiverParams{
		Logger: lc.logger, LogPrefix: lc.logPrefix,

		BrokerHost: commconf.BrokerServerHost(),

		Verbose: commconf.CommVerbose(),
	})
	if err != nil {
		return nil, err
	}

	_commBundle.streamSubscriber, err = streamPubsub.NewStreamSubscriber(&streamPubsub.StreamSubscriberParams{
		Context: lc.context,
		Logger: lc.logger, LogPrefix: lc.logPrefix,

		SnapshotServerHost: commconf.SnapshotServerHost(),
		PublisherServerHost: commconf.PublisherServerHost(),
	})
	if err != nil {
		return nil, err
	}

	return _commBundle, nil
}

func (commBundle *t_commBundle) start() (error) {
	var err error

	err = commBundle.tmpHubReceiver.Connect()
	if err != nil { return err }

	err = commBundle.nodeReceiver.Connect()
	if err != nil { return err }

	err = commBundle.nodeSubscriber.Connect()
	if err != nil { return err }

	err = commBundle.deviceReceiver.Connect()
	if err != nil { return err }

	err = commBundle.deviceSubscriber.Connect()
	if err != nil { return err }

	err = commBundle.streamReceiver.Connect()
	if err != nil { return err }

	err = commBundle.streamSubscriber.Connect()
	if err != nil { return err }

	return nil
}

func (commBundle *t_commBundle) stop() (error) {
	var err error

	err = commBundle.streamSubscriber.Disconnect()
	if err != nil { return err }

	err = commBundle.streamReceiver.Disconnect()
	if err != nil { return err }

	err = commBundle.deviceSubscriber.Disconnect()
	if err != nil { return err }

	err = commBundle.deviceReceiver.Disconnect()
	if err != nil { return err }

	err = commBundle.nodeSubscriber.Disconnect()
	if err != nil { return err }

	err = commBundle.nodeReceiver.Disconnect()
	if err != nil { return err }

	err = commBundle.tmpHubReceiver.Disconnect()
	if err != nil { return err }

	return nil
}
