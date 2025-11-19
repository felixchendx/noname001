package live

import (
	cacheIntface "noname001/app/module/common/cache/intface"
)

func (lc *LiveCache) makeCachedNode(_node *t_node) (*cacheIntface.CachedNode) {
	cachedNode := &cacheIntface.CachedNode{
		LastActivityAt: _node.lastExternalActivityAt,

		NodeSnapshot: *_node.nodeSnapshot,
		NodeStatus: lc.makeCachedNodeStatus(_node),
		NodeResource: lc.makeCachedNodeResource(_node),
		
		CachedDevices: lc.makeCachedDevices(_node),
		CachedStreams: lc.makeCachedStreams(_node),

		DevicesXStreams: make(map[string][]string),
	}

	for _, cachedDevice := range cachedNode.CachedDevices {
		streamCodeList := make([]string, 0)

		for _, cachedStream := range cachedNode.CachedStreams {
			if cachedStream.StreamSnapshot.Persistence.DeviceCode == cachedDevice.DeviceSnapshot.Persistence.Code {
				streamCodeList = append(streamCodeList, cachedStream.Code)
			}
		}

		cachedNode.DevicesXStreams[cachedDevice.DeviceSnapshot.Persistence.Code] = streamCodeList
	}

	return cachedNode
}

func (lc *LiveCache) makeCachedDevices(node *t_node) ([]*cacheIntface.CachedDevice) {
	var cachedDevices []*cacheIntface.CachedDevice

	cachedDevices = make([]*cacheIntface.CachedDevice, 0, len(node.deviceService.sortedDevices))

	for _, _device := range node.deviceService.sortedDevices {
		cachedDevices = append(cachedDevices, &cacheIntface.CachedDevice{
			LastUpdated   : _device.lastExternalActivityAt,
			DeviceSnapshot: *_device.deviceSnapshot,
			DeviceStatus  : lc.makeCachedDeviceStatus(node, _device),
		})
	}

	return cachedDevices
}

func (lc *LiveCache) makeCachedDevice(node *t_node, deviceCode string) (*cacheIntface.CachedDevice) {
	var cachedDevice *cacheIntface.CachedDevice

	seenDevice, alreadySeen := node.deviceService.devices[deviceCode]
	if alreadySeen {
		cachedDevice = &cacheIntface.CachedDevice{
			LastUpdated   : seenDevice.lastExternalActivityAt,
			DeviceSnapshot: *seenDevice.deviceSnapshot,
			DeviceStatus  : lc.makeCachedDeviceStatus(node, seenDevice),
		}
	}

	return cachedDevice
}

func (lc *LiveCache) makeCachedStreams(node *t_node) ([]*cacheIntface.CachedStream) {
	var cachedStreams []*cacheIntface.CachedStream

	cachedStreams = make([]*cacheIntface.CachedStream, 0, len(node.streamService.sortedStreams))

	for _, _stream := range node.streamService.sortedStreams {
		cachedStreams = append(cachedStreams, &cacheIntface.CachedStream{
			ID  : _stream.id,
			Code: _stream.code,

			StreamSnapshot: *_stream.streamSnapshot,
			StreamStatus: lc.makeCachedStreamStatus(node, _stream),

			SourceStream : _stream.sourceStream,
			RelayPathName: _stream.relayPathName,

			LastActivityAt: _stream.lastExternalActivityAt,
		})
	}

	return cachedStreams
}

func (lc *LiveCache) makeCachedStream(node *t_node, streamCode string) (*cacheIntface.CachedStream) {
	var cachedStream *cacheIntface.CachedStream

	seenStream, alreadySeen := node.streamService.streams[streamCode]
	if alreadySeen {
		cachedStream = &cacheIntface.CachedStream{
			ID  : seenStream.id,
			Code: seenStream.code,

			StreamSnapshot: *seenStream.streamSnapshot,
			StreamStatus: lc.makeCachedStreamStatus(node, seenStream),

			SourceStream : seenStream.sourceStream,
			RelayPathName: seenStream.relayPathName,

			LastActivityAt: seenStream.lastExternalActivityAt,
		}
	}

	return cachedStream
}
