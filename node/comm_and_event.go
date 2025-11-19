package node

import (
	nodeTyping "noname001/node/base/typing"
	nodeConst  "noname001/node/constant"
)

func (node *Node_t) _announce(evCode nodeTyping.NodeEventCode) {
	node.evHub.PublishNodeEvent(evCode, node.id)

	_ = node.commBundle.nodeEvPublisher.PublishNodeEvent(evCode, node.id, 60)
}

func (node *Node_t) setupHeartbeat() {
	node.cronJobs["heartbeat"], _ = node.cron.AddFunc(
		nodeConst.CROSSNODE_CRON_TIMING__NODE__HEARTBEAT,
		node._heartbeat,
	)
}

func (node *Node_t) _heartbeat() {
	// node.commBundle.nodeEvPublisher.PublishNodeLiveness(node.id, 60)
	// TODO: reconnect signal
	_ = node.commBundle.nodeEvPublisher.PublishNodeEvent(
		nodeTyping.NODE_EVENT_CODE__HEARTBEAT,
		node.id,
		60,
	)
}
