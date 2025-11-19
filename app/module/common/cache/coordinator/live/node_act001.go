package live

import (
	nodeTyping "noname001/node/base/typing"
)

func (lc *LiveCache) refreshNode(nodeID string) {
	nodeSnapshot, err := lc.fetchNodeSnapshot(nodeID)
	if err != nil {
		// activity log
		return
	}

	lc._nodeRefreshRoutine(nodeSnapshot)
}

func (lc *LiveCache) defunctNode(node *t_node, reason string) {
	lc.markNodeAsDefunct(node, reason)

	lc.defunctMediaServer(node, reason)
	lc.defunctDeviceService(node, reason)
	lc.defunctStreamService(node, reason)
}

func (lc *LiveCache) refreshNodeResource(node *t_node) {
	nodeResource, err := lc.fetchNodeResource(node.id)
	if err != nil {
		lc.updateNodeResource(node, nil)
		return
	}

	lc.updateNodeResource(node, nodeResource)
}

func (lc *LiveCache) _nodeRefreshRoutine(_nodeSnapshot *nodeTyping.BaseNodeSnapshot) {
	seenNode, alreadySeen := lc.nodes[_nodeSnapshot.ID]

	if alreadySeen {
		lc.updateNodeData(seenNode, _nodeSnapshot)

	} else {
		seenNode = lc.addNewNode(_nodeSnapshot)

		lc.refreshNodeResource(seenNode)
	}

	lc.interpretNodeStateToNodeStatus(seenNode)

	if seenNode.nodeSnapshot.AppSnapshot == nil {
		lc.defunctMediaServer(seenNode, "app:off")
		lc.defunctDeviceService(seenNode, "app:off")
		lc.defunctStreamService(seenNode, "app:off")

	} else {
		_, hasMediaServer := seenNode.nodeSnapshot.AppSnapshot.ModuleStates["common_mediasrv"]
		if hasMediaServer {
			lc.refreshMediaServer(seenNode)

		} else {
			lc.defunctMediaServer(seenNode, "mod:off")
		}

		_, hasDeviceService := seenNode.nodeSnapshot.AppSnapshot.ModuleStates["common_device"]
		if hasDeviceService {
			lc.refreshDeviceService(seenNode, true)

		} else {
			lc.defunctDeviceService(seenNode, "mod:off")
		}

		_, hasStreamService := seenNode.nodeSnapshot.AppSnapshot.ModuleStates["feature_stream"]
		if hasStreamService {
			lc.refreshStreamService(seenNode, true)

		} else {
			lc.defunctStreamService(seenNode, "mod:off")
		}
	}
}
