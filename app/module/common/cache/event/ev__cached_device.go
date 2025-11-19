package event

import (
	"sync"

	"github.com/google/uuid"

	deviceTyping "noname001/app/base/typing/device"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	cached_device_event_subscription__extra_cap int = 8
)

type CachedDeviceEvent = cacheIntface.CachedDeviceEvent

type CachedDeviceEventSubscription struct {
	id string

	Channel chan CachedDeviceEvent
}

type t_cachedDeviceEventSource struct {
	sourceChan chan CachedDeviceEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*CachedDeviceEventSubscription
	subscriptionsIndex map[string]*CachedDeviceEventSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newCachedDeviceEventSource() (*t_cachedDeviceEventSource) {
	evSource := &t_cachedDeviceEventSource{}
	evSource.sourceChan = make(chan CachedDeviceEvent)

	evSource.subscriptions = make([]*CachedDeviceEventSubscription, 0, cached_device_event_subscription__extra_cap)
	evSource.subscriptionsIndex = make(map[string]*CachedDeviceEventSubscription)

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

func (evHub *EventHub) PublishCachedDeviceEvent(oriDeviceEvent *deviceTyping.LiveDeviceEvent) {
	evHub.cachedDeviceEventSource.sourceChan <- CachedDeviceEvent{
		OriginalDeviceEvent: oriDeviceEvent,
	}
}

func (evHub *EventHub) NewCachedDeviceEventSubscription() (*CachedDeviceEventSubscription) {
	evSub := &CachedDeviceEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan CachedDeviceEvent)

	evHub.cachedDeviceEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveCachedDeviceEventSubscription(evSub *CachedDeviceEventSubscription) {
	evHub.cachedDeviceEventSource._removeSubscription(evSub)
}

func (evHub *EventHub) removeCachedDeviceEventSubscriptionByID(subID string) {
	evSub, inSubscription := evHub.cachedDeviceEventSource.subscriptionsIndex[subID]
	if inSubscription {
		evHub.cachedDeviceEventSource._removeSubscription(evSub)
	}
}

func (evSource *t_cachedDeviceEventSource) _addSubscription(evSub *CachedDeviceEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)
	evSource.subscriptionsIndex[evSub.id] = evSub

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_cachedDeviceEventSource) _removeSubscription(evSub *CachedDeviceEventSubscription) {
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

		resizedSubs := make([]*CachedDeviceEventSubscription, 0, subLen - 1 + cached_device_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
		delete(evSource.subscriptionsIndex, evSub.id)
	}

	evSource.subscriptionsMutex.Unlock()
}
