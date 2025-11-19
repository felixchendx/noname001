package event

import (
	"sync"

	"github.com/google/uuid"
)

const (
	STREAM_ITEM_EVENT_CODE__CREATE StreamItemEventCode = "si:create"
	STREAM_ITEM_EVENT_CODE__UPDATE StreamItemEventCode = "si:update"
	STREAM_ITEM_EVENT_CODE__DELETE StreamItemEventCode = "si:delete"

	stream_item_event_subscription__extra_cap int = 1
)

type StreamItemEventCode string
type StreamItemEvent struct {
	EventCode StreamItemEventCode

	StreamItemID   string
	StreamItemCode string
}

type StreamItemEventSubscription struct {
	id string

	Channel chan StreamItemEvent
}

type t_streamItemEventSource struct {
	sourceChan chan StreamItemEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*StreamItemEventSubscription
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newStreamItemEventSource() (*t_streamItemEventSource) {
	evSource := &t_streamItemEventSource{}
	evSource.sourceChan = make(chan StreamItemEvent)

	evSource.subscriptions = make([]*StreamItemEventSubscription, 0, stream_item_event_subscription__extra_cap)

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

func (evHub *EventHub) PublishStreamItemEvent(evCode StreamItemEventCode, streamItemID, streamItemCode string) {
	evHub.streamItemEventSource.sourceChan <- StreamItemEvent{
		EventCode: evCode,

		StreamItemID  : streamItemID,
		StreamItemCode: streamItemCode,
	}
}

func (evHub *EventHub) NewStreamItemEventSubscription() (*StreamItemEventSubscription) {
	evSub := &StreamItemEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan StreamItemEvent)

	evHub.streamItemEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveStreamItemEventSubscription(evSub *StreamItemEventSubscription) {
	evHub.streamItemEventSource._removeSubscription(evSub)
}

func (evSource *t_streamItemEventSource) _addSubscription(evSub *StreamItemEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_streamItemEventSource) _removeSubscription(evSub *StreamItemEventSubscription) {
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

		resizedSubs := make([]*StreamItemEventSubscription, 0, subLen - 1 + stream_item_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
	}

	evSource.subscriptionsMutex.Unlock()
}
