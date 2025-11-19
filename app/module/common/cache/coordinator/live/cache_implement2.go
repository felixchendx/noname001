package live

import (
	cacheIntface "noname001/app/module/common/cache/intface"
)

func (lc *LiveCache) CachedNodeStatus(nodeID string) (*cacheIntface.CachedNodeStatus) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedNodeStatus(seenNode)
	}

	return nil
}

func (lc *LiveCache) CachedNodeResource(nodeID string) (*cacheIntface.CachedNodeResource) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedNodeResource(seenNode)
	}

	return nil
}


func (lc *LiveCache) TempCachedDevicesAll() ([]*cacheIntface.CachedDevice) {
	return lc.makeTempCachedDevicesAll()
}

func (lc *LiveCache) CachedDeviceStatus(nodeID, deviceCode string) (*cacheIntface.CachedDeviceStatus) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedDeviceStatusXXX(seenNode, deviceCode)
	}

	return nil
}


func (lc *LiveCache) TempCachedStreamsAll() ([]*cacheIntface.CachedStream) {
	return lc.makeTempCachedStreamsAll()
}

func (lc *LiveCache) CachedStreamStatus(nodeID, streamCode string) (*cacheIntface.CachedStreamStatus) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedStreamStatusXXX(seenNode, streamCode)
	}

	return nil
}


func (lc *LiveCache) SubscribeToCachedNodeStatusFeed() (*cacheIntface.CachedNodeStatusFeedSubscription) {
	return lc.newNodeStatusFeedSubscription()
}

func (lc *LiveCache) UnsubscribeFromCachedNodeStatusFeed(dfSub *cacheIntface.CachedNodeStatusFeedSubscription) {
	lc.removeNodeStatusFeedSubscription(dfSub)
}

func (lc *LiveCache) SubscribeToCachedNodeResourceFeed() (*cacheIntface.CachedNodeResourceFeedSubscription) {
	return lc.newNodeResourceFeedSubscription()
}

func (lc *LiveCache) UnsubscribeFromCachedNodeResourceFeed(dfSub *cacheIntface.CachedNodeResourceFeedSubscription) {
	lc.removeNodeResourceFeedSubscription(dfSub)
}

func (lc *LiveCache) SubscribeToCachedDeviceStatusFeed() (*cacheIntface.CachedDeviceStatusFeedSubscription) {
	return lc.newDeviceStatusFeedSubscription()
}

func (lc *LiveCache) UnsubscribeFromCachedDeviceStatusFeed(dfSub *cacheIntface.CachedDeviceStatusFeedSubscription) {
	lc.removeDeviceStatusFeedSubscription(dfSub)
}

func (lc *LiveCache) SubscribeToCachedStreamStatusFeed() (*cacheIntface.CachedStreamStatusFeedSubscription) {
	return lc.newStreamStatusFeedSubscription()
}

func (lc *LiveCache) UnsubscribeFromCachedStreamStatusFeed(dfSub *cacheIntface.CachedStreamStatusFeedSubscription) {
	lc.removeStreamStatusFeedSubscription(dfSub)
}
