package live

import (
	"time"

	nodeTyping   "noname001/node/base/typing"
	deviceTyping "noname001/app/base/typing/device"
	streamTyping "noname001/app/base/typing/stream"
)

// TODO: network call optimization, i.e. folding events

// dedicated worker for each node
// mitigating node-to-node network problems
// i.e. some nodes might have higher latency / low bandwidth / general connectivity issues
//      that slow down 'main worker' that will adds up to halting
func (lc *LiveCache) nodeWorker(node *t_node) {
	var (
		datafeedTicker = time.NewTicker(5 * time.Second)
	)

	defer func() {
		datafeedTicker.Stop()
	}()

	looper:
	for {
		select {
		case <- node.context.Done():
			break looper

		case <- datafeedTicker.C:
			lc._datafeedRoutine(node)

		case nodeEv := <- node.nodeEventForwardingChan:
			lc._digestNodeEvent(node, nodeEv)

		case deviceEv := <- node.deviceEventForwardingChan:
			lc._digestDeviceEvent(node, deviceEv)

		case streamEv := <- node.streamEventForwardingChan:
			lc._digestStreamEvent(node, streamEv)
		}
	}
}

func (lc *LiveCache) _datafeedRoutine(node *t_node) {
	lc.refreshNodeResource(node)
	lc.feedNodeResource(node)
}

func (lc *LiveCache) _digestNodeEvent(node *t_node, nodeEv *nodeTyping.NodeEvent) {
	node.logExternalActivity("ev", []string{string(nodeEv.EventCode)})

	switch nodeEv.EventCode {
	case nodeTyping.NODE_EVENT_CODE__READY    : lc.refreshNode(node.id)
	case nodeTyping.NODE_EVENT_CODE__IP_CHANGE: lc.refreshNode(node.id)
	case nodeTyping.NODE_EVENT_CODE__SHUTDOWN : lc.defunctNode(node, string(nodeEv.EventCode))
	case nodeTyping.NODE_EVENT_CODE__HEARTBEAT: // noop
	default: // noop
	}

	lc.evHub.PublishCachedNodeEvent(nodeEv)

	lc.interpretNodeEventToNodeStatus(node, nodeEv)
	lc.feedNodeStatus(node)
}

func (lc *LiveCache) _digestDeviceEvent(node *t_node, deviceEv *deviceTyping.LiveDeviceEvent) {
	seenDevice, alreadySeen := node.deviceService.devices[deviceEv.DeviceCode]
	if alreadySeen {
		seenDevice.logExternalActivity("ev", []string{string(deviceEv.EventCode)})

		switch deviceEv.EventCode {
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__DEACTIVATED : lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_BEGIN  : lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_FAIL   : lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_OK     : lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__DISCONNECTED: lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_BEGIN: lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_FAIL : lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_OK   : lc.refreshDevice(node, deviceEv.DeviceCode)
		case deviceTyping.LIVE_DEVICE_EVENT_CODE__DESTROYED   : lc.defunctDevice(node, seenDevice, string(deviceEv.EventCode))
		default: // noop
		}

		lc.interpretDeviceEventToDeviceStatus(seenDevice, deviceEv)
		lc.feedDeviceStatus(node, seenDevice)

	} else {
		lc.refreshDevice(node, deviceEv.DeviceCode)
	}

	lc.evHub.PublishCachedDeviceEvent(deviceEv)
}

func (lc *LiveCache) _digestStreamEvent(node *t_node, streamEv *streamTyping.LiveStreamEvent) {
	seenStream, alreadySeen := node.streamService.streams[streamEv.StreamCode]
	if alreadySeen {
		seenStream.logExternalActivity("ev", []string{string(streamEv.EventCode)})

		switch streamEv.EventCode {
		case streamTyping.LIVE_STREAM_EVENT_CODE__DEACTIVATED : lc.refreshStream(node, streamEv.StreamCode)
		case streamTyping.LIVE_STREAM_EVENT_CODE__INIT_BEGIN  : lc.refreshStream(node, streamEv.StreamCode)
		case streamTyping.LIVE_STREAM_EVENT_CODE__INIT_FAIL   : lc.refreshStream(node, streamEv.StreamCode)
		case streamTyping.LIVE_STREAM_EVENT_CODE__INIT_OK     : lc.refreshStream(node, streamEv.StreamCode)
		case streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_BEGIN: lc.refreshStream(node, streamEv.StreamCode)
		case streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_FAIL : lc.refreshStream(node, streamEv.StreamCode)
		case streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_OK   : lc.refreshStream(node, streamEv.StreamCode)
		case streamTyping.LIVE_STREAM_EVENT_CODE__DESTROYED   : lc.defunctStream(node, seenStream, string(streamEv.EventCode))
		case streamTyping.LIVE_STREAM_EVENT_CODE__BG_FAIL     : lc.refreshStream(node, streamEv.StreamCode)
		default: // noop
		}

		lc.interpretStreamEventToStreamStatus(seenStream, streamEv)
		lc.feedStreamStatus(node, seenStream)

	} else {
		lc.refreshStream(node, streamEv.StreamCode)
	}

	lc.evHub.PublishCachedStreamEvent(streamEv)
}
