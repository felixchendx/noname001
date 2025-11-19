package live

// to be purged after these info can be accessed from web-gui

func (lc *LiveCache) dump() {
	lc.dumpInternal()
	// lc.dumpExternal()
}

func (lc *LiveCache) dumpInternal() {
	lc.logger.Debugf("===")
	lc.logger.Debugf("=== === === === === INTERNAL DUMP - BEGIN === === === === ===")

	for _, _node := range lc.sortedNodes {
		lc.logger.Debugf("=== === === === ===")
		lc.logger.Debugf("=== NODE - %s", _node.id)
		lc.logger.Debugf("===        currState: %s", _node.nodeSnapshot.State)
		lc.logger.Debugf("===        currIPs  : %s, AT: %s", _node.nodeSnapshot.IPs, _node.nodeSnapshot.LastIPHistoryTs)
		lc.logger.Debugf("===        localTime: %s, TZ: %s", _node.nodeSnapshot.LocalTime, _node.nodeSnapshot.Timezone)
		lc.logger.Debugf("===        modStates: %s", _node.nodeSnapshot.AppSnapshot.ModuleStates)
		// lc.logger.Debugf("===        int_logs[%v]: %s", len(_node.internalActivityLogs), _node.internalActivityLogs)
		// lc.logger.Debugf("===        ext_logs[%v]: %s", len(_node.externalActivityLogs), _node.externalActivityLogs)

		if _node.mediaServer != nil {
			lc.logger.Debugf("==")
			lc.logger.Debugf("==")
			lc.logger.Debugf("== MEDIA SERVER ==")
			lc.logger.Debugf("==              curr ip to use: %s", _node.mediaServer.ipToUse)
			lc.logger.Debugf("==              ports         : %s", _node.mediaServer.ports)
			lc.logger.Debugf("==              ping timestamp: %s", _node.mediaServer.pingTs)

			for _, _pingResult := range _node.mediaServer.pingResults {
				lc.logger.Debugf("==              _pingResult   : %s", _pingResult)
			}
		}

		if _node.deviceService != nil {
			lc.logger.Debugf("==")
			lc.logger.Debugf("==")
			lc.logger.Debugf("== DEVICE SERVICE ==")

			for _, _device := range _node.deviceService.sortedDevices {
				lc.logger.Debugf("==")
				lc.logger.Debugf("== DEVICE - %s, %s", _device.code, _device.id)
				lc.logger.Debugf("==          last activity at: %s", _device.lastExternalActivityAt)
				lc.logger.Debugf("==          snapshot: %s", _device.deviceSnapshot)
				// lc.logger.Debugf("==          int_logs[%v]: %s", len(_device.internalActivityLogs), _device.internalActivityLogs)
				// lc.logger.Debugf("==          ext_logs[%v]: %s", len(_device.externalActivityLogs), _device.externalActivityLogs)
			}
		}

		if _node.streamService != nil {
			lc.logger.Debugf("==")
			lc.logger.Debugf("==")
			lc.logger.Debugf("== STREAM SERVICE ==")
			
			for _, _stream := range _node.streamService.sortedStreams {
				lc.logger.Debugf("==")
				lc.logger.Debugf("== STREAM %s, %s", _stream.code, _stream.id)
				lc.logger.Debugf("==        %s", _stream.sourceStream)
				lc.logger.Debugf("==        %s", _stream.relayPathName)
				lc.logger.Debugf("==        %s", _stream.streamSnapshot)
			}
		}
	}

	lc.logger.Debugf("=== === === === === INTERNAL DUMP - ENDDD === === === === ===")
}

func (lc *LiveCache) dumpExternal() {
	lc.logger.Debugf("===")
	lc.logger.Debugf("=== === === === === CACHE DUMP - BEGIN === === === === ===")

	for _, _cachedNode := range lc.CachedNodes() {
		lc.logger.Debugf("=== === === === ===")
		lc.logger.Debugf("=== NODE CACHE - %s", _cachedNode.NodeSnapshot.ID)
		lc.logger.Debugf("===              device: %v, stream: %v", len(_cachedNode.CachedDevices), len(_cachedNode.CachedStreams))

		lc.logger.Debugf("==")
		lc.logger.Debugf("== DEVICE CACHES ==")
		for _, _cachedDevice := range _cachedNode.CachedDevices {
			lc.logger.Debugf("== DEVICE CACHE - %s, %6s", _cachedDevice.DeviceSnapshot.Persistence.Code, _cachedDevice.DeviceSnapshot.Persistence.State)
			lc.logger.Debugf("==                %s", _cachedDevice)
		}

		lc.logger.Debugf("==")
		lc.logger.Debugf("== STREAM CACHES ==")
		for _, _cachedStream := range _cachedNode.CachedStreams {
			lc.logger.Debugf("==")
			lc.logger.Debugf("== STREAM CACHE - %s", _cachedStream.Code)
			lc.logger.Debugf("==                %s", _cachedStream.RelayPathName)
			lc.logger.Debugf("==                %s", _cachedStream.SourceStream)
			lc.logger.Debugf("==                %s", _cachedStream.StreamSnapshot)
		}

		lc.logger.Debugf("==")
		lc.logger.Debugf("== DEVICES X STREAMS ==")
		for _deviceCode, _streamCodeList := range _cachedNode.DevicesXStreams {
			lc.logger.Debugf("== %12s - %s", _deviceCode, _streamCodeList)
		}
	}

	lc.logger.Debugf("=== === === === === CACHE DUMP - ENDDD === === === === ===")
}
