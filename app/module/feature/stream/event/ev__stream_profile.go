package event

import (
	"sync"

	"github.com/google/uuid"
)

const (
	STREAM_PROFILE_EVENT_CODE__CREATE StreamProfileEventCode = "sp:create"
	STREAM_PROFILE_EVENT_CODE__UPDATE StreamProfileEventCode = "sp:update"
	STREAM_PROFILE_EVENT_CODE__DELETE StreamProfileEventCode = "sp:delete"

	stream_profile_event_subscription__extra_cap int = 1
)

type StreamProfileEventCode string
type StreamProfileEvent struct {
	EventCode StreamProfileEventCode

	StreamProfileID   string
	StreamProfileCode string
}

type StreamProfileEventSubscription struct {
	id string

	Channel chan StreamProfileEvent
}

type t_streamProfileEventSource struct {
	sourceChan chan StreamProfileEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*StreamProfileEventSubscription
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newStreamProfileEventSource() (*t_streamProfileEventSource) {
	evSource := &t_streamProfileEventSource{}
	evSource.sourceChan = make(chan StreamProfileEvent)

	evSource.subscriptions = make([]*StreamProfileEventSubscription, 0, stream_profile_event_subscription__extra_cap)

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

func (evHub *EventHub) PublishStreamProfileEvent(evCode StreamProfileEventCode, streamProfileID, streamProfileCode string) {
	evHub.streamProfileEventSource.sourceChan <- StreamProfileEvent{
		EventCode: evCode,

		StreamProfileID  : streamProfileID,
		StreamProfileCode: streamProfileCode,
	}
}

func (evHub *EventHub) NewStreamProfileEventSubscription() (*StreamProfileEventSubscription) {
	evSub := &StreamProfileEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan StreamProfileEvent)

	evHub.streamProfileEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveStreamProfileEventSubscription(evSub *StreamProfileEventSubscription) {
	evHub.streamProfileEventSource._removeSubscription(evSub)
}

func (evSource *t_streamProfileEventSource) _addSubscription(evSub *StreamProfileEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_streamProfileEventSource) _removeSubscription(evSub *StreamProfileEventSubscription) {
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

		resizedSubs := make([]*StreamProfileEventSubscription, 0, subLen - 1 + stream_profile_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
	}

	evSource.subscriptionsMutex.Unlock()
}
