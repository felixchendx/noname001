package live

import (
	"sync"

	"github.com/google/uuid"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	node_resource_feed_subscription__extra_cap int = 8
)

type CachedNodeResource = cacheIntface.CachedNodeResource
type CachedNodeResourceFeedSubscription = cacheIntface.CachedNodeResourceFeedSubscription

type t_nodeResourceFeedSource struct {
	sourceChan chan *CachedNodeResource

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*CachedNodeResourceFeedSubscription
	subscriptionsIndex map[string]*CachedNodeResourceFeedSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (lc *LiveCache) initNodeResourceFeedSource() {
	dfSource := &t_nodeResourceFeedSource{}
	dfSource.sourceChan = make(chan *CachedNodeResource)

	dfSource.subscriptions = make([]*CachedNodeResourceFeedSubscription, 0, node_resource_feed_subscription__extra_cap)
	dfSource.subscriptionsIndex = make(map[string]*CachedNodeResourceFeedSubscription)

	go func() {
		dfSourceLooper:
		for {
			select {
			case <- lc.context.Done():
				break dfSourceLooper

			case dat := <- dfSource.sourceChan:
				for _, dfSub := range dfSource.subscriptions {
					dfSub.Channel <- dat
				}
			}
		}
	}()

	lc.nodeResourceFeedSource = dfSource
}

func (lc *LiveCache) feedNodeResource(node *t_node) {
	lc.nodeResourceFeedSource.sourceChan <- lc.makeCachedNodeResource(node)
}

func (lc *LiveCache) newNodeResourceFeedSubscription() (*CachedNodeResourceFeedSubscription) {
	dfSub := &CachedNodeResourceFeedSubscription{}
	dfSub.ID = uuid.New().String()
	dfSub.Channel = make(chan *CachedNodeResource)

	lc.nodeResourceFeedSource._addSubscription(dfSub)

	return dfSub
}

func (lc *LiveCache) removeNodeResourceFeedSubscription(dfSub *CachedNodeResourceFeedSubscription) {
	lc.nodeResourceFeedSource._removeSubscription(dfSub)
}

func (dfSource *t_nodeResourceFeedSource) _addSubscription(dfSub *CachedNodeResourceFeedSubscription) {
	dfSource.subscriptionsMutex.Lock()
	defer dfSource.subscriptionsMutex.Unlock()

	dfSource.subscriptions = append(dfSource.subscriptions, dfSub)
	dfSource.subscriptionsIndex[dfSub.ID] = dfSub
}

func (dfSource *t_nodeResourceFeedSource) _removeSubscription(dfSub *CachedNodeResourceFeedSubscription) {
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

		resizedSubs := make([]*CachedNodeResourceFeedSubscription, 0, subLen - 1 + node_resource_feed_subscription__extra_cap)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[subIdx+1:subLen]...)

		dfSource.subscriptions = resizedSubs
		delete(dfSource.subscriptionsIndex, dfSub.ID)
	}

	// instead of dynamic dynamic buffer management
	// try high and low watermark indicator
}
