package event

import (
	baseEv "noname001/internal/base/event"
)

type EventTree string
const (
	EVTREE__BACKGROUND string = "/local/bg/"
)

type EventCode string
const (
)

type BgEvent = baseEv.BaseBgEvent

func (evHub *EventHub) NewBgEvent(action string) (*BgEvent) {
	return baseEv.NewBaseBgEvent(action)
}

func (evHub *EventHub) PublishBgEvent(ev *BgEvent) {
	evHub.PublishToEventSource(EVTREE__BACKGROUND, ev)
}
