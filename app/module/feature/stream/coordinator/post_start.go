package coordinator

import (
	localEv "noname001/dilemma/event"
)

func (coord *Coordinator) postStartRoutine() {
	coord.postStart__loadLiveStream()
}

func (coord *Coordinator) postStart__loadLiveStream() {
	bgev := localEv.LocalEventHub().NewBgEvent("postStart__loadLiveStream")

	streamItemList, dbev := coord.store.DB.Coord__postStart_loadLiveStream()
	if dbev.IsError() {
		bgev.Messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return
	}

	for _, streamItemPE := range streamItemList {
		coord.initializeLiveStream(streamItemPE.ID)
	}
}
