package event

import (
	localEv "noname001/dilemma/event"
)

type EventHub struct {
	*localEv.EventHub
}

func ExtendLocalEventHub() (*EventHub) {
	evHub := &EventHub{localEv.LocalEventHub()}
	return evHub
}

func (evHub *EventHub) Open() {
}

func (evHub *EventHub) Close() {
}
