package event

import (
	"sync"
	"time"

	"github.com/google/uuid"

	deviceTyping "noname001/app/base/typing/device"
)

const (
	live_device_event_subscription__extra_cap int = 2
)

type LiveDeviceEventCode = deviceTyping.LiveDeviceEventCode
type LiveDeviceEvent     = deviceTyping.LiveDeviceEvent

type LiveDeviceEventSubscription struct {
	id string

	Channel chan LiveDeviceEvent
}

type t_liveDeviceEventSource struct {
	sourceChan chan LiveDeviceEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*LiveDeviceEventSubscription
	subscriptionsIndex map[string]*LiveDeviceEventSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newLiveDeviceEventSource() (*t_liveDeviceEventSource) {
	evSource := &t_liveDeviceEventSource{}
	evSource.sourceChan = make(chan LiveDeviceEvent)

	evSource.subscriptions = make([]*LiveDeviceEventSubscription, 0, live_device_event_subscription__extra_cap)
	evSource.subscriptionsIndex = make(map[string]*LiveDeviceEventSubscription)

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

func (evHub *EventHub) PublishLiveDeviceEvent(evCode LiveDeviceEventCode, nodeID, deviceID, deviceCode string) {
	evHub.liveDeviceEventSource.sourceChan <- LiveDeviceEvent{
		Timestamp: time.Now(),
		EventCode: evCode,

		NodeID    : nodeID,
		DeviceID  : deviceID,
		DeviceCode: deviceCode,
	}
}

func (evHub *EventHub) NewLiveDeviceEventSubscription() (*LiveDeviceEventSubscription) {
	evSub := &LiveDeviceEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan LiveDeviceEvent)

	evHub.liveDeviceEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveLiveDeviceEventSubscription(evSub *LiveDeviceEventSubscription) {
	evHub.liveDeviceEventSource._removeSubscription(evSub)
}

func (evHub *EventHub) removeLiveDeviceEventSubscriptionByID(subID string) {
	evSub, inSubscription := evHub.liveDeviceEventSource.subscriptionsIndex[subID]
	if inSubscription {
		evHub.liveDeviceEventSource._removeSubscription(evSub)
	}
}

func (evSource *t_liveDeviceEventSource) _addSubscription(evSub *LiveDeviceEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)
	evSource.subscriptionsIndex[evSub.id] = evSub

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_liveDeviceEventSource) _removeSubscription(evSub *LiveDeviceEventSubscription) {
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

		resizedSubs := make([]*LiveDeviceEventSubscription, 0, subLen - 1 + live_device_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
		delete(evSource.subscriptionsIndex, evSub.id)
	}

	evSource.subscriptionsMutex.Unlock()
}
