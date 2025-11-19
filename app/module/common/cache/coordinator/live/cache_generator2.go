package live

import (
	cacheIntface "noname001/app/module/common/cache/intface"
)

func (lc *LiveCache) makeCachedNodeStatus(_node *t_node) (*cacheIntface.CachedNodeStatus) {
	var cachedNodeStatus = &cacheIntface.CachedNodeStatus{}
	cachedNodeStatus.ID = _node.id

	if _node.nodeStatus != nil {
		cachedNodeStatus.Timestamp = _node.nodeStatus.timestamp

		cachedNodeStatus.TextualIndicator = _node.nodeStatus.textualIndicator
		cachedNodeStatus.VisualIndicator = _node.nodeStatus.visualIndicator
		cachedNodeStatus.AuditoryIndicator = _node.nodeStatus.auditoryIndicator
	}

	return cachedNodeStatus
}

func (lc *LiveCache) makeCachedNodeResource(_node *t_node) (*cacheIntface.CachedNodeResource) {
	if _node.nodeResource == nil { return nil }

	cachedNodeResource := &cacheIntface.CachedNodeResource{
		ID: _node.id,

		NodeResource: *_node.nodeResource,
	}

	return cachedNodeResource
}

func (lc *LiveCache) makeTempCachedDevicesAll() ([]*cacheIntface.CachedDevice) {
	var (
		allCachedDevicesCount int = 0
		allCachedDevices      []*cacheIntface.CachedDevice
		cachedDevicesByNodes  [][]*cacheIntface.CachedDevice
	)
	cachedDevicesByNodes = make([][]*cacheIntface.CachedDevice, 0, len(lc.sortedNodes))

	for _, _node := range lc.sortedNodes {
		var cachedDevicesForCurrNode = lc.makeCachedDevices(_node)

		cachedDevicesByNodes = append(cachedDevicesByNodes, cachedDevicesForCurrNode)
		allCachedDevicesCount += len(cachedDevicesForCurrNode)
	}

	allCachedDevices = make([]*cacheIntface.CachedDevice, 0, allCachedDevicesCount)
	for _, _slice := range cachedDevicesByNodes {
		allCachedDevices = append(allCachedDevices, _slice...)
	}

	return allCachedDevices
}

func (lc *LiveCache) makeCachedDeviceStatus(node *t_node, device *t_device) (*cacheIntface.CachedDeviceStatus) {
	var cachedDeviceStatus = &cacheIntface.CachedDeviceStatus{}
	cachedDeviceStatus.DeviceID = device.id

	cachedDeviceStatus.NodeID = node.id
	cachedDeviceStatus.DeviceCode = device.code

	if device.deviceStatus != nil {
		cachedDeviceStatus.Timestamp = device.deviceStatus.timestamp

		cachedDeviceStatus.TextualIndicator = device.deviceStatus.textualIndicator
		cachedDeviceStatus.VisualIndicator = device.deviceStatus.visualIndicator
		cachedDeviceStatus.AuditoryIndicator = device.deviceStatus.auditoryIndicator
	}

	return cachedDeviceStatus
}

// TODO
func (lc *LiveCache) makeCachedDeviceStatusXXX(node *t_node, deviceCode string) (*cacheIntface.CachedDeviceStatus) {
	seenDevice, alreadySeen := node.deviceService.devices[deviceCode]
	if alreadySeen {
		return lc.makeCachedDeviceStatus(node, seenDevice)
	}

	var cachedDeviceStatus = &cacheIntface.CachedDeviceStatus{}
	cachedDeviceStatus.NodeID = node.id
	cachedDeviceStatus.DeviceCode = deviceCode

	return cachedDeviceStatus
}

func (lc *LiveCache) makeTempCachedStreamsAll() ([]*cacheIntface.CachedStream) {
	var (
		allCachedStreamsCount int = 0
		allCachedStreams      []*cacheIntface.CachedStream
		cachedStreamsByNodes  [][]*cacheIntface.CachedStream
	)
	cachedStreamsByNodes = make([][]*cacheIntface.CachedStream, 0, len(lc.sortedNodes))

	for _, _node := range lc.sortedNodes {
		var cachedStreamsForCurrNode = lc.makeCachedStreams(_node)

		cachedStreamsByNodes = append(cachedStreamsByNodes, cachedStreamsForCurrNode)
		allCachedStreamsCount += len(cachedStreamsForCurrNode)
	}

	allCachedStreams = make([]*cacheIntface.CachedStream, 0, allCachedStreamsCount)
	for _, _slice := range cachedStreamsByNodes {
		allCachedStreams = append(allCachedStreams, _slice...)
	}

	return allCachedStreams
}

func (lc *LiveCache) makeCachedStreamStatus(node *t_node, stream *t_stream) (*cacheIntface.CachedStreamStatus) {
	var cachedStreamStatus = &cacheIntface.CachedStreamStatus{}
	cachedStreamStatus.StreamID = stream.id

	cachedStreamStatus.NodeID = node.id
	cachedStreamStatus.StreamCode = stream.code

	if stream.streamStatus != nil {
		cachedStreamStatus.Timestamp = stream.streamStatus.timestamp

		cachedStreamStatus.TextualIndicator = stream.streamStatus.textualIndicator
		cachedStreamStatus.VisualIndicator = stream.streamStatus.visualIndicator
		cachedStreamStatus.AuditoryIndicator = stream.streamStatus.auditoryIndicator
	}

	return cachedStreamStatus
}

// TODO
func (lc *LiveCache) makeCachedStreamStatusXXX(node *t_node, streamCode string) (*cacheIntface.CachedStreamStatus) {
	seenStream, alreadySeen := node.streamService.streams[streamCode]
	if alreadySeen {
		return lc.makeCachedStreamStatus(node, seenStream)
	}

	var cachedStreamStatus = &cacheIntface.CachedStreamStatus{}
	cachedStreamStatus.NodeID = node.id
	cachedStreamStatus.StreamCode = streamCode

	return cachedStreamStatus
}
