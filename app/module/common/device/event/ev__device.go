package event

import (
	"sync"

	"github.com/google/uuid"
)

const (
	DEVICE_EVENT_CODE__CREATE DeviceEventCode = "d:create"
	DEVICE_EVENT_CODE__UPDATE DeviceEventCode = "d:update"
	DEVICE_EVENT_CODE__DELETE DeviceEventCode = "d:delete"

	device_event_subscription__extra_cap int = 1
)

type DeviceEventCode string
type DeviceEvent struct {
	// EventID   string
	EventCode DeviceEventCode

	DeviceID   string
	DeviceCode string
}

type DeviceEventSubscription struct {
	id string

	Channel chan DeviceEvent
}

type t_deviceEventSource struct {
	sourceChan chan DeviceEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*DeviceEventSubscription
	// === subscriptionsMutex ===
}

func (evHub *EventHub) newDeviceEventSource() (*t_deviceEventSource) {
	evSource := &t_deviceEventSource{}
	evSource.sourceChan = make(chan DeviceEvent)

	evSource.subscriptions = make([]*DeviceEventSubscription, 0, 8)

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

func (evHub *EventHub) PublishDeviceEvent(evCode DeviceEventCode, deviceID, deviceCode string) {
	evHub.deviceEventSource.sourceChan <- DeviceEvent{
		EventCode: evCode,

		DeviceID  : deviceID,
		DeviceCode: deviceCode,
	}
}

func (evHub *EventHub) NewDeviceEventSubscription() (*DeviceEventSubscription) {
	evSub := &DeviceEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan DeviceEvent)

	evHub.deviceEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveDeviceEventSubscription(evSub *DeviceEventSubscription) {
	evHub.deviceEventSource._removeSubscription(evSub)
}

func (evSource *t_deviceEventSource) _addSubscription(evSub *DeviceEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_deviceEventSource) _removeSubscription(evSub *DeviceEventSubscription) {
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

		resizedSubs := make([]*DeviceEventSubscription, 0, subLen - 1 + device_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
	}

	evSource.subscriptionsMutex.Unlock()
}
