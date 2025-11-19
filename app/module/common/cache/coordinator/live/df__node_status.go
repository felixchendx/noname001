package live

import (
	"sync"

	"github.com/google/uuid"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	node_status_feed_subscription__extra_cap int = 8
)

type CachedNodeStatus = cacheIntface.CachedNodeStatus
type CachedNodeStatusFeedSubscription = cacheIntface.CachedNodeStatusFeedSubscription

type t_nodeStatusFeedSource struct {
	sourceChan chan *CachedNodeStatus

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*CachedNodeStatusFeedSubscription
	subscriptionsIndex map[string]*CachedNodeStatusFeedSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (lc *LiveCache) initNodeStatusFeedSource() {
	dfSource := &t_nodeStatusFeedSource{}
	dfSource.sourceChan = make(chan *CachedNodeStatus)

	dfSource.subscriptions = make([]*CachedNodeStatusFeedSubscription, 0, node_status_feed_subscription__extra_cap)
	dfSource.subscriptionsIndex = make(map[string]*CachedNodeStatusFeedSubscription)

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

	lc.nodeStatusFeedSource = dfSource
}

func (lc *LiveCache) feedNodeStatus(node *t_node) {
	lc.nodeStatusFeedSource.sourceChan <- lc.makeCachedNodeStatus(node)
}

func (lc *LiveCache) newNodeStatusFeedSubscription() (*CachedNodeStatusFeedSubscription) {
	dfSub := &CachedNodeStatusFeedSubscription{}
	dfSub.ID = uuid.New().String()
	dfSub.Channel = make(chan *CachedNodeStatus)

	lc.nodeStatusFeedSource._addSubscription(dfSub)

	return dfSub
}

func (lc *LiveCache) removeNodeStatusFeedSubscription(dfSub *CachedNodeStatusFeedSubscription) {
	lc.nodeStatusFeedSource._removeSubscription(dfSub)
}

func (dfSource *t_nodeStatusFeedSource) _addSubscription(dfSub *CachedNodeStatusFeedSubscription) {
	dfSource.subscriptionsMutex.Lock()
	defer dfSource.subscriptionsMutex.Unlock()

	dfSource.subscriptions = append(dfSource.subscriptions, dfSub)
	dfSource.subscriptionsIndex[dfSub.ID] = dfSub
}

func (dfSource *t_nodeStatusFeedSource) _removeSubscription(dfSub *CachedNodeStatusFeedSubscription) {
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

		resizedSubs := make([]*CachedNodeStatusFeedSubscription, 0, subLen - 1 + node_status_feed_subscription__extra_cap)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[subIdx+1:subLen]...)

		dfSource.subscriptions = resizedSubs
		delete(dfSource.subscriptionsIndex, dfSub.ID)
	}

	// instead of dynamic dynamic buffer management
	// try high and low watermark indicator
}
