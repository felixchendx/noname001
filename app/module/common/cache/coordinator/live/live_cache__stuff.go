package live

import (
	"strings"
	"sync"
	"time"
)

var (
	zero_time time.Time
)

func (lc *LiveCache) scanRoutine() {
	nodeInfoList, err := lc.fetchNodeInfoList()
	if err != nil {
		// scan activity log
		return
	}

	wg := new(sync.WaitGroup)
	for _, nodeInfo := range nodeInfoList {
		wg.Add(1)
		go func() {
			lc.refreshNode(nodeInfo.NodeID) // space out each call abit ?
			wg.Done()
		}()
	}
	wg.Wait()
}

func (lc *LiveCache) cleanupRoutine() {
	var checkTime = time.Now()

	for _, node := range lc.sortedNodes {
		lc.checkNodeStaleness(node, checkTime)
		nodeExpired := lc.isNodeExpired(node, checkTime)

		if nodeExpired {
			lc.removeNode(node)

		} else {
			for _, device := range node.deviceService.sortedDevices {
				node.deviceService.checkDeviceStaleness(device, checkTime)
				deviceExpired := node.deviceService.isDeviceExpired(device, checkTime)
				if deviceExpired {
					node.deviceService.removeDevice(device)
				}
			}

			for _, stream := range node.streamService.sortedStreams {
				node.streamService.checkStreamStaleness(stream, checkTime)
				streamExpired := node.streamService.isStreamExpired(stream, checkTime)
				if streamExpired {
					node.streamService.removeStream(stream)
				}
			}
		}
	}
}

const (
	CACHE_STATE__OK    CacheState = "c:ok"
	CACHE_STATE__EOL   CacheState = "c:end_of_life"
	CACHE_STATE__STALE CacheState = "c:stale"
)

type CacheState string

type t_internalActivityLog struct {
	ts time.Time
	activity string
	result string
}

type t_externalActivityLog struct {
	ts time.Time
	activity string
	extra []string
}

func _caseInsensitiveSort(a, b string) (int) {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b))
}
