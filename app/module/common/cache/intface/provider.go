package intface

type CacheDataProviderIntface interface {
	// TODO: manual controls / overrides

	CachedNodes()             ([]*CachedNode)
	CachedNode(nodeID string) (*CachedNode)

	CachedDevices(nodeID string)            ([]*CachedDevice)
	CachedDevice(nodeID, deviceCode string) (*CachedDevice)

	CachedStreams(nodeID string)            ([]*CachedStream)
	CachedStream(nodeID, streamCode string) (*CachedStream)

	// === 2 temp ? these stuffs might be moved to mod_monitoring ===
	CachedNodeStatus(nodeID string)   (*CachedNodeStatus)
	CachedNodeResource(nodeID string) (*CachedNodeResource)

	TempCachedDevicesAll()                        ([]*CachedDevice)
	CachedDeviceStatus(nodeID, deviceCode string) (*CachedDeviceStatus)

	TempCachedStreamsAll()                        ([]*CachedStream)
	CachedStreamStatus(nodeID, streamCode string) (*CachedStreamStatus)

	SubscribeToCachedNodeStatusFeed() (*CachedNodeStatusFeedSubscription)
	UnsubscribeFromCachedNodeStatusFeed(*CachedNodeStatusFeedSubscription)
	SubscribeToCachedNodeResourceFeed() (*CachedNodeResourceFeedSubscription)
	UnsubscribeFromCachedNodeResourceFeed(*CachedNodeResourceFeedSubscription)

	SubscribeToCachedDeviceStatusFeed() (*CachedDeviceStatusFeedSubscription)
	UnsubscribeFromCachedDeviceStatusFeed(*CachedDeviceStatusFeedSubscription)

	SubscribeToCachedStreamStatusFeed() (*CachedStreamStatusFeedSubscription)
	UnsubscribeFromCachedStreamStatusFeed(*CachedStreamStatusFeedSubscription)
}

type CacheEventProviderIntface interface {
	SubscribeToCachedNodeEvent() (*CachedNodeEventSubscription)
	UnsubscribeFromCachedNodeEvent(*CachedNodeEventSubscription)

	SubscribeToCachedDeviceEvent() (*CachedDeviceEventSubscription)
	UnsubscribeFromCachedDeviceEvent(*CachedDeviceEventSubscription)

	SubscribeToCachedStreamEvent() (*CachedStreamEventSubscription)
	UnsubscribeFromCachedStreamEvent(*CachedStreamEventSubscription)
}

func DataProvider() (CacheDataProviderIntface) {
	return cacheDataProvider
}

func EventProvider() (CacheEventProviderIntface) {
	return cacheEventProvider
}

// === ^^^ for those that uses     ^^^ ===
// =======================================
// === VVV for those that provides VVV ===

var (
	cacheDataProvider  CacheDataProviderIntface
	cacheEventProvider CacheEventProviderIntface
)

func AssignCacheDataProvider(_something CacheDataProviderIntface) {
	cacheDataProvider = _something
}

func AssignCacheEventProvider(_something CacheEventProviderIntface) {
	cacheEventProvider = _something
}
