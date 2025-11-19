package coordinator

import (
	"noname001/app/module/feature/stream/event"
)

// only instance registry related, i.e. initialize and terminate
func (coord *Coordinator) streamEventListeners() {
	var (
		streamGroupEvSub = coord.evHub.NewStreamGroupEventSubscription()
		streamItemEvSub  = coord.evHub.NewStreamItemEventSubscription()
	)

	defer func() {
		coord.evHub.RemoveStreamGroupEventSubscription(streamGroupEvSub)
		coord.evHub.RemoveStreamItemEventSubscription(streamItemEvSub)
	}()

	evListenerLoop:
	for {
		select {
		case <- coord.context.Done():
			break evListenerLoop

		case ev := <- streamGroupEvSub.Channel:
			switch ev.EventCode {
			case event.STREAM_GROUP_EVENT_CODE__CREATE: // noop
			case event.STREAM_GROUP_EVENT_CODE__UPDATE: // noop, handled inside live stream
			case event.STREAM_GROUP_EVENT_CODE__DELETE: coord.terminateLiveStreamByGroup(ev.StreamGroupID)
			default:
			}

		case ev := <- streamItemEvSub.Channel:
			switch ev.EventCode {
			case event.STREAM_ITEM_EVENT_CODE__CREATE: coord.initializeLiveStream(ev.StreamItemID)
			case event.STREAM_ITEM_EVENT_CODE__UPDATE: // noop, handled inside live stream
			case event.STREAM_ITEM_EVENT_CODE__DELETE: coord.terminateLiveStream(ev.StreamItemID)
			default:
			}
		}
	}
}

// temp, grouping is to be dismissed for unknown period of time
func (coord *Coordinator) terminateLiveStreamByGroup(groupID string) {
	for k, liveStream := range coord.liveStreams {
		inGroup := liveStream.DestroyIfBelongToThisGroup(groupID)
		if inGroup {
			delete(coord.liveStreams, k)
		}
	}
}
