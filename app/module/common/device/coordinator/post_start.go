package coordinator

import (
	localEv "noname001/dilemma/event"
)

func (coord *Coordinator) postStartRoutine() {
	go coord.loadAutoStartDevice()
}

func (coord *Coordinator) loadAutoStartDevice() {
	bgev := localEv.LocalEventHub().NewBgEvent("loadAutoStartDevice")

	deviceIDList, dbev := coord.store.DB.Coord__loadAutoStartDevice()
	if dbev.IsError() {
		bgev.Messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		bgev.Messages.AddError(POST_ERR_01001.NewMessage())
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return
	}

	for _, id := range deviceIDList {
		coord.initLiveDevice(id)
	}
}
