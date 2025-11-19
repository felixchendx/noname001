package event

import (
	streamIntface "noname001/app/module/feature/stream/intface"
)

func (evHub *EventHub) SubscribeToLiveStreamEvent() (*streamIntface.LiveStreamEventSubscription) {
	var evSub           = evHub.NewLiveStreamEventSubscription()
	var interfacedEvSub = &streamIntface.LiveStreamEventSubscription{
		ID     : evSub.id,
		Channel: evSub.Channel,
	}

	return interfacedEvSub
}

func (evHub *EventHub) UnsubscribeFromLiveStreamEvent(interfacedEvSub *streamIntface.LiveStreamEventSubscription) {
	evHub.removeLiveStreamEventSubscriptionByID(interfacedEvSub.ID)
}
