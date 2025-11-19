package node

import (
	"time"

	"noname001/app/base/messaging"

	nodeTyping "noname001/node/base/typing"
	nodeMdp    "noname001/node/comm/mdpservice"
	nodePubsub "noname001/node/comm/pubsub"

	mediasrvIntface "noname001/app/module/common/mediasrv/intface"

	nodePSUtil "noname001/node/psutil"
)

type nodeCommBundle_t struct {
	// TODO: node -> hub connection state, and persistent reconnect
	//       might need to replace zmq examples with internal implementation
	nodeProvider    *nodeMdp.NodeProvider
	nodeEvPublisher *nodePubsub.NodePublisher
}

func (node *Node_t) initComm() (*nodeCommBundle_t, error) {
	var (
		bundle *nodeCommBundle_t = &nodeCommBundle_t{}
		err    error
	)

	bundle.nodeProvider, err = nodeMdp.NewNodeProvider(&nodeMdp.NodeProviderParams{
		Context: node.context,
		Logger: node.logger, LogPrefix: node.logPrefix,

		TempNodeID : node.id,
		BrokerHost : node.cfgRoot.Node.BrokerServerHost,
		DataHandler: node,

		Verbose     : node.cfgRoot.Node.CommVerbose,
		RetryBackoff: nodeMdp.DEFAULT_RETRY_BACKOFF,
	})
	if err != nil {
		node.logger.Errorf("%s: initComm err %s", node.logPrefix, err.Error())
		return nil, err
	}

	bundle.nodeEvPublisher, err = nodePubsub.NewNodePublisher(&nodePubsub.NodePublisherParams{
		Context: node.context,
		Logger: node.logger, LogPrefix: node.logPrefix,

		CollectorServerHost: node.cfgRoot.Node.CollectorServerHost,
	})
	if err != nil {
		node.logger.Errorf("%s: initComm err %s", node.logPrefix, err.Error())
		return nil, err
	}

	return bundle, nil
}

func (commBundle *nodeCommBundle_t) connect() (error) {
	var err error

	err = commBundle.nodeProvider.Connect()
	if err != nil {
		return err
	}

	err = commBundle.nodeEvPublisher.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (commBundle *nodeCommBundle_t) disconnect() {
	_ = commBundle.nodeEvPublisher.Disconnect()
	_ = commBundle.nodeProvider.Disconnect()
}

// ============ VVV conform to nodeMdp.DataHandlerIntface VVV ============= //
func (node *Node_t) ProvideNodeSnapshot() (*nodeTyping.BaseNodeSnapshot, *messaging.Messages) {
	messages := messaging.NewMessages()

	nodeSnapshot := &nodeTyping.BaseNodeSnapshot{
		ID   : node.id,
		Name : node.name,
		State: node.state,

		LocalTime: time.Now(),
		Timezone:  node.cfgRoot.Global.Timezone,

		IPs: node.ips,
	}

	if len(node.ipCollectionHistories) > 0 {
		lastIPHistory := node.ipCollectionHistories[len(node.ipCollectionHistories)-1]
		nodeSnapshot.LastIPHistoryTs = lastIPHistory.timestamp
	}

	if node.hub != nil {
		nodeSnapshot.HubSnapshot = &nodeTyping.BaseHubSnapshot{}
	}

	if node.app != nil {
		loadedModuleStates := node.app.ModStates()

		nodeSnapshot.AppSnapshot = &nodeTyping.BaseAppSnapshot{
			ModuleStates: loadedModuleStates,
		}

		_, hasMediaServer := loadedModuleStates["common_mediasrv"]
		if hasMediaServer {
			mediasrvProvider := mediasrvIntface.Provider()
			nodeSnapshot.AppSnapshot.TempMediasrvSnapshot = &nodeTyping.TempMediasrvSnapshot{
				Ports    : mediasrvProvider.StreamingPorts(),
				AuthnPair: mediasrvProvider.RelayAuthnPair(),
			}
		}
	}

	return nodeSnapshot, messages
}

func (node *Node_t) ProvideNodeResource() (*nodeTyping.TempNodeSystemResourceSummary, *messaging.Messages) {
	messages := messaging.NewMessages()
	tempSumm := nodePSUtil.TempSystemResourceSummary()

	return tempSumm, messages
}
// ============ ^^^ conform to nodeMdp.DataHandlerIntface ^^^ ============= //
