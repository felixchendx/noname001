package event

import (
	"sync"

	"github.com/google/uuid"
)

const (
	STREAM_GROUP_EVENT_CODE__CREATE StreamGroupEventCode = "sg:create"
	STREAM_GROUP_EVENT_CODE__UPDATE StreamGroupEventCode = "sg:update"
	STREAM_GROUP_EVENT_CODE__DELETE StreamGroupEventCode = "sg:delete"

	stream_group_event_subscription__extra_cap int = 1
)

type StreamGroupEventCode string
type StreamGroupEvent struct {
	EventCode StreamGroupEventCode

	StreamGroupID   string
	StreamGroupCode string
}

type StreamGroupEventSubscription struct {
	id string

	Channel chan StreamGroupEvent
}

type t_streamGroupEventSource struct {
	sourceChan chan StreamGroupEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*StreamGroupEventSubscription
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newStreamGroupEventSource() (*t_streamGroupEventSource) {
	evSource := &t_streamGroupEventSource{}
	evSource.sourceChan = make(chan StreamGroupEvent)

	evSource.subscriptions = make([]*StreamGroupEventSubscription, 0, stream_group_event_subscription__extra_cap)

	go func() {
		evSourceLoop:
		for {
			select {
			case <- evHub.context.Done():
				break evSourceLoop

			case ev := <- evSource.sourceChan:
				for _, evSub := range evSource.subscriptions {
					evSub.Channel <- ev
				}
			}
		}
	}()

	return evSource
}

func (evHub *EventHub) PublishStreamGroupEvent(evCode StreamGroupEventCode, streamGroupID, streamGroupCode string) {
	evHub.streamGroupEventSource.sourceChan <- StreamGroupEvent{
		EventCode: evCode,

		StreamGroupID  : streamGroupID,
		StreamGroupCode: streamGroupCode,
	}
}

func (evHub *EventHub) NewStreamGroupEventSubscription() (*StreamGroupEventSubscription) {
	evSub := &StreamGroupEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan StreamGroupEvent)

	evHub.streamGroupEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveStreamGroupEventSubscription(evSub *StreamGroupEventSubscription) {
	evHub.streamGroupEventSource._removeSubscription(evSub)
}

func (evSource *t_streamGroupEventSource) _addSubscription(evSub *StreamGroupEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_streamGroupEventSource) _removeSubscription(evSub *StreamGroupEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	subIdx := -1

	for _subIdx, _evSub := range evSource.subscriptions {
		if _evSub.id == evSub.id {
			subIdx = _subIdx
			break
		}
	}

	if subIdx != -1 {
		subLen := len(evSource.subscriptions)

		resizedSubs := make([]*StreamGroupEventSubscription, 0, subLen - 1 + stream_group_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
	}

	evSource.subscriptionsMutex.Unlock()
}
