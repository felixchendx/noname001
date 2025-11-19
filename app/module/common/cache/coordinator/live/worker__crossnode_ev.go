package live

import (
	nodeTyping   "noname001/node/base/typing"
	deviceTyping "noname001/app/base/typing/device"
	streamTyping "noname001/app/base/typing/stream"
)

// no blocking calls on this worker
func (lc *LiveCache) crossnodeEventListeners() {
	var (
		nodeEvSub   = lc.commBundle.nodeSubscriber.SubscribeToNodeEvent()
		deviceEvSub = lc.commBundle.deviceSubscriber.SubscribeToDeviceEvent()
		streamEvSub = lc.commBundle.streamSubscriber.SubscribeToStreamEvent()
	)

	defer func() {
		lc.commBundle.nodeSubscriber.UnsubscribeFromNodeEvent(nodeEvSub)
		lc.commBundle.deviceSubscriber.UnsubscribeFromDeviceEvent(deviceEvSub)
		lc.commBundle.streamSubscriber.UnsubscribeFromStreamEvent(streamEvSub)
	}()

	evListenerLoop:
	for {
		select {
		case <- lc.context.Done():
			break evListenerLoop

		case ev := <- nodeEvSub.MessageChannel:
			lc._nodeEventForwarding(ev)

		case ev := <- deviceEvSub.MessageChannel:
			lc._deviceEventForwarding(ev)

		case ev := <- streamEvSub.MessageChannel:
			lc._streamEventForwarding(ev)
		}
	}
}

func (lc *LiveCache) _nodeEventForwarding(ev *nodeTyping.NodeEvent) {
	seenNode, alreadySeen := lc.nodes[ev.NodeID]
	if alreadySeen {
		seenNode.nodeEventForwardingChan <- ev
	} else {
		lc.miscJobChan <- &t_miscJob{misc_job__node_not_seen, []string{ev.NodeID}}
	}
}

func (lc *LiveCache) _deviceEventForwarding(ev *deviceTyping.LiveDeviceEvent) {
	seenNode, alreadySeen := lc.nodes[ev.NodeID]
	if alreadySeen {
		seenNode.deviceEventForwardingChan <- ev
	} else {
		lc.miscJobChan <- &t_miscJob{misc_job__node_not_seen, []string{ev.NodeID}}
	}
}

func (lc *LiveCache) _streamEventForwarding(ev *streamTyping.LiveStreamEvent) {
	seenNode, alreadySeen := lc.nodes[ev.NodeID]
	if alreadySeen {
		seenNode.streamEventForwardingChan <- ev
	} else {
		lc.miscJobChan <- &t_miscJob{misc_job__node_not_seen, []string{ev.NodeID}}
	}
}
