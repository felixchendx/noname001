package event

import (
	"sync"

	"github.com/google/uuid"

	nodeTyping "noname001/node/base/typing"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	cached_node_event_subscription__extra_cap int = 8
)

type CachedNodeEvent = cacheIntface.CachedNodeEvent

type CachedNodeEventSubscription struct {
	id string

	Channel chan CachedNodeEvent
}

type t_cachedNodeEventSource struct {
	sourceChan chan CachedNodeEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*CachedNodeEventSubscription
	subscriptionsIndex map[string]*CachedNodeEventSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newCachedNodeEventSource() (*t_cachedNodeEventSource) {
	evSource := &t_cachedNodeEventSource{}
	evSource.sourceChan = make(chan CachedNodeEvent)

	evSource.subscriptions = make([]*CachedNodeEventSubscription, 0, cached_node_event_subscription__extra_cap)
	evSource.subscriptionsIndex = make(map[string]*CachedNodeEventSubscription)

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

func (evHub *EventHub) PublishCachedNodeEvent(originalNodeEvent *nodeTyping.NodeEvent) {
	evHub.cachedNodeEventSource.sourceChan <- CachedNodeEvent{
		OriginalNodeEvent: originalNodeEvent,
	}
}

func (evHub *EventHub) NewCachedNodeEventSubscription() (*CachedNodeEventSubscription) {
	evSub := &CachedNodeEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan CachedNodeEvent)

	evHub.cachedNodeEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveCachedNodeEventSubscription(evSub *CachedNodeEventSubscription) {
	evHub.cachedNodeEventSource._removeSubscription(evSub)
}

func (evHub *EventHub) removeCachedNodeEventSubscriptionByID(subID string) {
	evSub, inSubscription := evHub.cachedNodeEventSource.subscriptionsIndex[subID]
	if inSubscription {
		evHub.cachedNodeEventSource._removeSubscription(evSub)
	}
}

func (evSource *t_cachedNodeEventSource) _addSubscription(evSub *CachedNodeEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)
	evSource.subscriptionsIndex[evSub.id] = evSub

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_cachedNodeEventSource) _removeSubscription(evSub *CachedNodeEventSubscription) {
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

		resizedSubs := make([]*CachedNodeEventSubscription, 0, subLen - 1 + cached_node_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
		delete(evSource.subscriptionsIndex, evSub.id)
	}

	evSource.subscriptionsMutex.Unlock()
}
