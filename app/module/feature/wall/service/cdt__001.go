package service

import (
	cacheIntface "noname001/app/module/common/cache/intface"
)

func (svc *Service) GetCachedNodes() ([]*cacheIntface.CachedNode) {
	return svc.coordinator.GetCachedNodes()
}

func (svc *Service) GetCachedNode(nodeID string) (*cacheIntface.CachedNode) {
	return svc.coordinator.GetCachedNode(nodeID)
}

func (svc *Service) GetCachedStreams(nodeID string) ([]*cacheIntface.CachedStream) {
	return svc.coordinator.GetCachedStreams(nodeID)
}

func (svc *Service) GetCachedStream(nodeID, streamCode string) (*cacheIntface.CachedStream) {
	return svc.coordinator.GetCachedStream(nodeID, streamCode)
}

func (svc *Service) GetRelayedStreamViewURL(requesterHostname, nodeID, streamCode, streamProtocol string) (string) {
	return svc.coordinator.GetRelayedStreamViewURL(requesterHostname, nodeID, streamCode, streamProtocol)
}
