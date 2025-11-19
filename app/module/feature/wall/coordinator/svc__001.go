package coordinator

import (
	mediasrvIntface "noname001/app/module/common/mediasrv/intface"
	cacheIntface    "noname001/app/module/common/cache/intface"
)

func (coord *Coordinator) GetCachedNodes() ([]*cacheIntface.CachedNode) {
	return cacheIntface.DataProvider().CachedNodes()
}

func (coord *Coordinator) GetCachedNode(nodeID string) (*cacheIntface.CachedNode) {
	return cacheIntface.DataProvider().CachedNode(nodeID)
}

func (coord *Coordinator) GetCachedStreams(nodeID string) ([]*cacheIntface.CachedStream) {
	return cacheIntface.DataProvider().CachedStreams(nodeID)
}

func (coord *Coordinator) GetCachedStream(nodeID, streamCode string) (*cacheIntface.CachedStream) {
	return cacheIntface.DataProvider().CachedStream(nodeID, streamCode)
}

func (coord *Coordinator) GetRelayedStreamViewURL(requesterHostname, nodeID, streamCode, streamProtocol string) (string) {
	var streamURL string

	cachedStream := cacheIntface.DataProvider().CachedStream(nodeID, streamCode)
	if cachedStream != nil {
		streamURL = mediasrvIntface.Provider().RelayedStreamViewURL(requesterHostname, nodeID, streamCode, streamProtocol)
	}

	return streamURL
}
