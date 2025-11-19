package event

import (
	"sync"

	"github.com/google/uuid"

	streamTyping "noname001/app/base/typing/stream"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	cached_stream_event_subscription__extra_cap int = 8
)

type CachedStreamEvent = cacheIntface.CachedStreamEvent

type CachedStreamEventSubscription struct {
	id string

	Channel chan CachedStreamEvent
}

type t_cachedStreamEventSource struct {
	sourceChan chan CachedStreamEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*CachedStreamEventSubscription
	subscriptionsIndex map[string]*CachedStreamEventSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newCachedStreamEventSource() (*t_cachedStreamEventSource) {
	evSource := &t_cachedStreamEventSource{}
	evSource.sourceChan = make(chan CachedStreamEvent)

	evSource.subscriptions = make([]*CachedStreamEventSubscription, 0, cached_stream_event_subscription__extra_cap)
	evSource.subscriptionsIndex = make(map[string]*CachedStreamEventSubscription)

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

func (evHub *EventHub) PublishCachedStreamEvent(oriStreamEvent *streamTyping.LiveStreamEvent) {
	evHub.cachedStreamEventSource.sourceChan <- CachedStreamEvent{
		OriginalStreamEvent: oriStreamEvent,
	}
}

func (evHub *EventHub) NewCachedStreamEventSubscription() (*CachedStreamEventSubscription) {
	evSub := &CachedStreamEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan CachedStreamEvent)

	evHub.cachedStreamEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveCachedStreamEventSubscription(evSub *CachedStreamEventSubscription) {
	evHub.cachedStreamEventSource._removeSubscription(evSub)
}

func (evHub *EventHub) removeCachedStreamEventSubscriptionByID(subID string) {
	evSub, inSubscription := evHub.cachedStreamEventSource.subscriptionsIndex[subID]
	if inSubscription {
		evHub.cachedStreamEventSource._removeSubscription(evSub)
	}
}

func (evSource *t_cachedStreamEventSource) _addSubscription(evSub *CachedStreamEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)
	evSource.subscriptionsIndex[evSub.id] = evSub

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_cachedStreamEventSource) _removeSubscription(evSub *CachedStreamEventSubscription) {
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

		resizedSubs := make([]*CachedStreamEventSubscription, 0, subLen - 1 + cached_stream_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
		delete(evSource.subscriptionsIndex, evSub.id)
	}

	evSource.subscriptionsMutex.Unlock()
}
