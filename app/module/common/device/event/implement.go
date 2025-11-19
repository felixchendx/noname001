package event

import (
	deviceIntface "noname001/app/module/common/device/intface"
)

func (evHub *EventHub) SubscribeToLiveDeviceEvent() (*deviceIntface.LiveDeviceEventSubscription) {
	var evSub           = evHub.NewLiveDeviceEventSubscription()
	var interfacedEvSub = &deviceIntface.LiveDeviceEventSubscription{
		ID     : evSub.id,
		Channel: evSub.Channel,
	}

	return interfacedEvSub
}

func (evHub *EventHub) UnsubscribeFromLiveDeviceEvent(interfacedEvSub *deviceIntface.LiveDeviceEventSubscription) {
	evHub.removeLiveDeviceEventSubscriptionByID(interfacedEvSub.ID)
}
