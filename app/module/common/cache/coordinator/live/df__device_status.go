package live

import (
	"sync"

	"github.com/google/uuid"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	device_status_feed_subscription__extra_cap int = 8
)

type CachedDeviceStatus = cacheIntface.CachedDeviceStatus
type CachedDeviceStatusFeedSubscription = cacheIntface.CachedDeviceStatusFeedSubscription

type t_deviceStatusFeedSource struct {
	sourceChan chan *CachedDeviceStatus

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*CachedDeviceStatusFeedSubscription
	subscriptionsIndex map[string]*CachedDeviceStatusFeedSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (lc *LiveCache) initDeviceStatusFeedSource() {
	dfSource := &t_deviceStatusFeedSource{}
	dfSource.sourceChan = make(chan *CachedDeviceStatus)

	dfSource.subscriptions = make([]*CachedDeviceStatusFeedSubscription, 0, device_status_feed_subscription__extra_cap)
	dfSource.subscriptionsIndex = make(map[string]*CachedDeviceStatusFeedSubscription)

	go func() {
		dfSourceLoop:
		for {
			select {
			case <- lc.context.Done():
				break dfSourceLoop

			case dat := <- dfSource.sourceChan:
				for _, dfSub := range dfSource.subscriptions {
					dfSub.Channel <- dat
				}
			}
		}
	}()

	lc.deviceStatusFeedSource = dfSource
}

func (lc *LiveCache) feedDeviceStatus(node *t_node, device *t_device) {
	lc.deviceStatusFeedSource.sourceChan <- lc.makeCachedDeviceStatus(node, device)
}

func (lc *LiveCache) newDeviceStatusFeedSubscription() (*CachedDeviceStatusFeedSubscription) {
	dfSub := &CachedDeviceStatusFeedSubscription{}
	dfSub.ID = uuid.New().String()
	dfSub.Channel = make(chan *CachedDeviceStatus)

	lc.deviceStatusFeedSource._addSubscription(dfSub)

	return dfSub
}

func (lc *LiveCache) removeDeviceStatusFeedSubscription(dfSub *CachedDeviceStatusFeedSubscription) {
	lc.deviceStatusFeedSource._removeSubscription(dfSub)
}

func (dfSource *t_deviceStatusFeedSource) _addSubscription(dfSub *CachedDeviceStatusFeedSubscription) {
	dfSource.subscriptionsMutex.Lock()
	defer dfSource.subscriptionsMutex.Unlock()

	dfSource.subscriptions = append(dfSource.subscriptions, dfSub)
	dfSource.subscriptionsIndex[dfSub.ID] = dfSub
}

func (dfSource *t_deviceStatusFeedSource) _removeSubscription(dfSub *CachedDeviceStatusFeedSubscription) {
	dfSource.subscriptionsMutex.Lock()
	defer dfSource.subscriptionsMutex.Unlock()

	subIdx := -1

	for _subIdx, _dfSub := range dfSource.subscriptions {
		if _dfSub.ID == dfSub.ID {
			subIdx = _subIdx
			break
		}
	}

	if subIdx != -1 {
		subLen := len(dfSource.subscriptions)

		resizedSubs := make([]*CachedDeviceStatusFeedSubscription, 0, subLen - 1 + device_status_feed_subscription__extra_cap)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[subIdx+1:subLen]...)

		dfSource.subscriptions = resizedSubs
		delete(dfSource.subscriptionsIndex, dfSub.ID)
	}

	// instead of dynamic dynamic buffer management
	// try high and low watermark indicator
}
