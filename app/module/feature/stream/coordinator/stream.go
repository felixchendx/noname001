package coordinator

import (
	"noname001/app/module/feature/stream/coordinator/live"
)

func (coord *Coordinator) initializeLiveStream(streamID string) {
	existingLiveStream, alreadyExist := coord.liveStreams[streamID]
	if alreadyExist {
		existingLiveStream.Init() // is idempotent

	} else {
		liveStream := live.NewLiveStream(&live.LiveStreamParams{
			ParentContext: coord.context,
			ParentLogger : coord.logger, ParentLogPrefix: coord.logPrefix,

			EvHub        : coord.evHub,
			Store        : coord.store,
			CommBundle   : coord.commBundle,

			StreamID     : streamID,
		})
	
		coord.liveStreams[streamID] = liveStream
		liveStream.Init()
	}
}

func (coord *Coordinator) reloadLiveStream(streamID string) {
	existingLiveStream, alreadyExist := coord.liveStreams[streamID]
	if !alreadyExist {
		return
	}

	existingLiveStream.Reload()
}

func (coord *Coordinator) terminateLiveStream(streamID string) {
	existingLiveStream, alreadyExist := coord.liveStreams[streamID]
	if !alreadyExist {
		return
	}

	existingLiveStream.Destroy()
	delete(coord.liveStreams, streamID)
}
