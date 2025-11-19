package live

import (
	streamEv "noname001/app/module/feature/stream/event"

	deviceTyping  "noname001/app/base/typing/device"
	deviceIntface "noname001/app/module/common/device/intface"
)

func (liveStream *LiveStream) evListener() {
	var (
		liveDeviceEvSub = deviceIntface.EventProvider().SubscribeToLiveDeviceEvent()

		streamProfileEvSub = liveStream.evHub.NewStreamProfileEventSubscription()
		streamGroupEvSub   = liveStream.evHub.NewStreamGroupEventSubscription()
		streamItemEvSub    = liveStream.evHub.NewStreamItemEventSubscription()
	)

	defer func() {
		deviceIntface.EventProvider().UnsubscribeFromLiveDeviceEvent(liveDeviceEvSub)

		liveStream.evHub.RemoveStreamProfileEventSubscription(streamProfileEvSub)
		liveStream.evHub.RemoveStreamGroupEventSubscription(streamGroupEvSub)
		liveStream.evHub.RemoveStreamItemEventSubscription(streamItemEvSub)
	}()

	evListenerLoop:
	for {
		selectCase:
		select {
		case <- liveStream.context.Done():
			break evListenerLoop

		case _ev := <- liveDeviceEvSub.Channel:
			if liveStream.pDat.SourceType != "mod_device" { break selectCase }
			if liveStream.pDat.DeviceCode != _ev.DeviceCode { break selectCase }

			// TODO: dep oriented exec
			switch _ev.EventCode {
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__DEACTIVATED : liveStream.execChan <- exec_reload
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_BEGIN  : // noop
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_FAIL   : liveStream.execChan <- exec_reload
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_OK     : liveStream.execChan <- exec_reload
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__DISCONNECTED: liveStream.execChan <- exec_reload
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_BEGIN: // noop
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_FAIL : liveStream.execChan <- exec_reload
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_OK   : liveStream.execChan <- exec_reload
			case deviceTyping.LIVE_DEVICE_EVENT_CODE__DESTROYED   : liveStream.execChan <- exec_reload
			default:
			}

		case _ev := <- streamProfileEvSub.Channel:
			if liveStream.pDat.ProfileID != _ev.StreamProfileID { break selectCase }

			switch _ev.EventCode {
			case streamEv.STREAM_PROFILE_EVENT_CODE__CREATE: // noop
			case streamEv.STREAM_PROFILE_EVENT_CODE__UPDATE: liveStream.execChan <- exec_reload
			case streamEv.STREAM_PROFILE_EVENT_CODE__DELETE: // noop, guarded by service
			default:
			}

		case _ev := <- streamGroupEvSub.Channel:
			if liveStream.pDat.GroupID != _ev.StreamGroupID { break selectCase }

			switch _ev.EventCode {
			case streamEv.STREAM_GROUP_EVENT_CODE__CREATE: // noop
			case streamEv.STREAM_GROUP_EVENT_CODE__UPDATE: liveStream.execChan <- exec_reload
			case streamEv.STREAM_GROUP_EVENT_CODE__DELETE: // noop, governed by coordinator
			default:
			}

		case _ev := <- streamItemEvSub.Channel:
			if liveStream.id != _ev.StreamItemID { break selectCase }

			switch _ev.EventCode {
			case streamEv.STREAM_ITEM_EVENT_CODE__CREATE: // noop, governed by coordinator
			case streamEv.STREAM_ITEM_EVENT_CODE__UPDATE: liveStream.execChan <- exec_reload
			case streamEv.STREAM_ITEM_EVENT_CODE__DELETE: // noop, governed by coordinator
			default:
			}
		}
	}
}
