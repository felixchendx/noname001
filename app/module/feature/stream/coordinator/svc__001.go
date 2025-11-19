package coordinator

import (
	mediasrvIntface "noname001/app/module/common/mediasrv/intface"

	liveStream "noname001/app/module/feature/stream/coordinator/live"
)

func (coord *Coordinator) GetLiveStreams() (map[string]*liveStream.LiveStream) {
	return coord.liveStreams
}

func (coord *Coordinator) GetLiveStreamByCode(streamCode string) (*liveStream.LiveStream, bool) {
	for _, _liveStream := range coord.liveStreams {
		if _liveStream.Code() == streamCode {
			return _liveStream, true
		}
	}

	return nil, false
}

func (coord *Coordinator) GetStreamViewURL(requesterHostname, streamItemCode, streamProtocol string) (string) {
	return mediasrvIntface.Provider().StreamViewURL(requesterHostname, streamItemCode, streamProtocol)
}
