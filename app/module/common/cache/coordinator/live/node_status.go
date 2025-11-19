package live

import (
	"time"

	nodeTyping "noname001/node/base/typing"
)

func (lc *LiveCache) interpretNodeEventToNodeStatus(node *t_node, nodeEv *nodeTyping.NodeEvent) {
	var nodeStatus = &t_nodeStatusInterpretation{
		timestamp: time.Now(),
	}

	switch nodeEv.EventCode {
	case nodeTyping.NODE_EVENT_CODE__READY:
		nodeStatus.textualIndicator = "ready"
		nodeStatus.visualIndicator  = visual_indicator__green_steady

	case nodeTyping.NODE_EVENT_CODE__IP_CHANGE:
		nodeStatus.textualIndicator = "ip_change"
		nodeStatus.visualIndicator  = visual_indicator__yellow_blink

	case nodeTyping.NODE_EVENT_CODE__SHUTDOWN:
		nodeStatus.textualIndicator = "shutdown"
		nodeStatus.visualIndicator  = visual_indicator__off

	case nodeTyping.NODE_EVENT_CODE__HEARTBEAT:
		nodeStatus.textualIndicator = "heartbeat"
		nodeStatus.visualIndicator  = visual_indicator__green_steady
		// nodeStatus.visualIndicator  = visual_indicator_sequence__twest

	default:
		// noop
	}

	node.nodeStatus = nodeStatus
}

func (lc *LiveCache) interpretNodeStateToNodeStatus(node *t_node) {
	var nodeStatus = &t_nodeStatusInterpretation{
		timestamp: time.Now(),
	}

	switch node.nodeSnapshot.State {
	case nodeTyping.NODE_STATE__INIT:
		nodeStatus.textualIndicator = "initializing"
		nodeStatus.visualIndicator  = visual_indicator__green_blink

	case nodeTyping.NODE_STATE__START:
		nodeStatus.textualIndicator = "starting"
		nodeStatus.visualIndicator  = visual_indicator__green_blink
		
	case nodeTyping.NODE_STATE__READY:
		nodeStatus.textualIndicator = "ready"
		nodeStatus.visualIndicator  = visual_indicator__green_steady

	case nodeTyping.NODE_STATE__ABORT:
		nodeStatus.textualIndicator = "abort"
		nodeStatus.visualIndicator  = visual_indicator__red_blink
		
	case nodeTyping.NODE_STATE__STOP:
		fallthrough
	case nodeTyping.NODE_STATE__SHUTDOWN:
		nodeStatus.textualIndicator = "shutdown"
		nodeStatus.visualIndicator  = visual_indicator__off

	case nodeTyping.NODE_STATE__DISCONNECTED:
		nodeStatus.textualIndicator = "shutdown"
		nodeStatus.visualIndicator  = visual_indicator__yellow_blink

	default:
		// noop
	}

	node.nodeStatus = nodeStatus
}
