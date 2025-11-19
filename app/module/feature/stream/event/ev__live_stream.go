package event

import (
	"sync"
	"time"

	"github.com/google/uuid"

	streamTyping "noname001/app/base/typing/stream"
)

const (
	live_stream_event_subscription__extra_cap int = 2
)

type LiveStreamEventCode = streamTyping.LiveStreamEventCode
type LiveStreamEvent     = streamTyping.LiveStreamEvent

type LiveStreamEventSubscription struct {
	id string

	Channel chan LiveStreamEvent
}

type t_liveStreamEventSource struct {
	sourceChan chan LiveStreamEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*LiveStreamEventSubscription
	subscriptionsIndex map[string]*LiveStreamEventSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newLiveStreamEventSource() (*t_liveStreamEventSource) {
	evSource := &t_liveStreamEventSource{}
	evSource.sourceChan = make(chan LiveStreamEvent)

	evSource.subscriptions = make([]*LiveStreamEventSubscription, 0, live_stream_event_subscription__extra_cap)
	evSource.subscriptionsIndex = make(map[string]*LiveStreamEventSubscription)

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

func (evHub *EventHub) PublishLiveStreamEvent(evCode LiveStreamEventCode, nodeID, streamID, streamCode string) {
	evHub.liveStreamEventSource.sourceChan <- LiveStreamEvent{
		Timestamp: time.Now(),
		EventCode: evCode,

		NodeID    : nodeID,
		StreamID  : streamID,
		StreamCode: streamCode,
	}
}

func (evHub *EventHub) NewLiveStreamEventSubscription() (*LiveStreamEventSubscription) {
	evSub := &LiveStreamEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan LiveStreamEvent)

	evHub.liveStreamEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveLiveStreamEventSubscription(evSub *LiveStreamEventSubscription) {
	evHub.liveStreamEventSource._removeSubscription(evSub)
}

func (evHub *EventHub) removeLiveStreamEventSubscriptionByID(subID string) {
	evSub, inSubscription := evHub.liveStreamEventSource.subscriptionsIndex[subID]
	if inSubscription {
		evHub.liveStreamEventSource._removeSubscription(evSub)
	}
}

func (evSource *t_liveStreamEventSource) _addSubscription(evSub *LiveStreamEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)
	evSource.subscriptionsIndex[evSub.id] = evSub

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_liveStreamEventSource) _removeSubscription(evSub *LiveStreamEventSubscription) {
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

		resizedSubs := make([]*LiveStreamEventSubscription, 0, subLen - 1 + live_stream_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
		delete(evSource.subscriptionsIndex, evSub.id)
	}

	evSource.subscriptionsMutex.Unlock()
}
