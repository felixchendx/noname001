package event

import (
	cacheIntface "noname001/app/module/common/cache/intface"
)

func (evHub *EventHub) SubscribeToCachedNodeEvent() (*cacheIntface.CachedNodeEventSubscription) {
	var evSub           = evHub.NewCachedNodeEventSubscription()
	var interfacedEvSub = &cacheIntface.CachedNodeEventSubscription{
		ID     : evSub.id,
		Channel: evSub.Channel,
	}

	return interfacedEvSub
}

func (evHub *EventHub) UnsubscribeFromCachedNodeEvent(interfacedEvSub *cacheIntface.CachedNodeEventSubscription) {
	evHub.removeCachedNodeEventSubscriptionByID(interfacedEvSub.ID)
}

func (evHub *EventHub) SubscribeToCachedDeviceEvent() (*cacheIntface.CachedDeviceEventSubscription) {
	var evSub           = evHub.NewCachedDeviceEventSubscription()
	var interfacedEvSub = &cacheIntface.CachedDeviceEventSubscription{
		ID     : evSub.id,
		Channel: evSub.Channel,
	}

	return interfacedEvSub
}

func (evHub *EventHub) UnsubscribeFromCachedDeviceEvent(interfacedEvSub *cacheIntface.CachedDeviceEventSubscription) {
	evHub.removeCachedDeviceEventSubscriptionByID(interfacedEvSub.ID)
}

func (evHub *EventHub) SubscribeToCachedStreamEvent() (*cacheIntface.CachedStreamEventSubscription) {
	var evSub           = evHub.NewCachedStreamEventSubscription()
	var interfacedEvSub = &cacheIntface.CachedStreamEventSubscription{
		ID     : evSub.id,
		Channel: evSub.Channel,
	}

	return interfacedEvSub
}

func (evHub *EventHub) UnsubscribeFromCachedStreamEvent(interfacedEvSub *cacheIntface.CachedStreamEventSubscription) {
	evHub.removeCachedStreamEventSubscriptionByID(interfacedEvSub.ID)
}
