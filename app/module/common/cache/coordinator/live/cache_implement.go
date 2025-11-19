package live

import (
	cacheIntface "noname001/app/module/common/cache/intface"
)

// no pre-cache for now, too much additional complexity

func (lc *LiveCache) CachedNodes() ([]*cacheIntface.CachedNode) {
	var nodesCache = make([]*cacheIntface.CachedNode, len(lc.sortedNodes))

	for nodeIdx, node := range lc.sortedNodes {
		nodesCache[nodeIdx] = lc.makeCachedNode(node)
	}

	return nodesCache
}

func (lc *LiveCache) CachedNode(nodeID string) (*cacheIntface.CachedNode) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedNode(seenNode)
	}

	return nil
}

func (lc *LiveCache) CachedDevices(nodeID string) ([]*cacheIntface.CachedDevice) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedDevices(seenNode)
	}

	return make([]*cacheIntface.CachedDevice, 0)
}

func (lc *LiveCache) CachedDevice(nodeID, deviceCode string) (*cacheIntface.CachedDevice) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedDevice(seenNode, deviceCode)
	}

	return nil
}

func (lc *LiveCache) CachedStreams(nodeID string) ([]*cacheIntface.CachedStream) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedStreams(seenNode)
	}

	return make([]*cacheIntface.CachedStream, 0)
}

func (lc *LiveCache) CachedStream(nodeID, streamCode string) (*cacheIntface.CachedStream) {
	seenNode, alreadySeen := lc.nodes[nodeID]
	if alreadySeen {
		return lc.makeCachedStream(seenNode, streamCode)
	}

	return nil
}
