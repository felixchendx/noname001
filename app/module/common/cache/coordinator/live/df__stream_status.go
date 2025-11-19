package live

import (
	"sync"

	"github.com/google/uuid"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	stream_status_feed_subscription__extra_cap int = 8
)

type CachedStreamStatus = cacheIntface.CachedStreamStatus
type CachedStreamStatusFeedSubscription = cacheIntface.CachedStreamStatusFeedSubscription

type t_streamStatusFeedSource struct {
	sourceChan chan *CachedStreamStatus

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*CachedStreamStatusFeedSubscription
	subscriptionsIndex map[string]*CachedStreamStatusFeedSubscription // k: subscriptionID
	// === subscriptionsMutex ===
}

func (lc *LiveCache) initStreamStatusFeedSource() {
	dfSource := &t_streamStatusFeedSource{}
	dfSource.sourceChan = make(chan *CachedStreamStatus)

	dfSource.subscriptions = make([]*CachedStreamStatusFeedSubscription, 0, stream_status_feed_subscription__extra_cap)
	dfSource.subscriptionsIndex = make(map[string]*CachedStreamStatusFeedSubscription)

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

	lc.streamStatusFeedSource = dfSource
}

func (lc *LiveCache) feedStreamStatus(node *t_node, device *t_stream) {
	lc.streamStatusFeedSource.sourceChan <- lc.makeCachedStreamStatus(node, device)
}

func (lc *LiveCache) newStreamStatusFeedSubscription() (*CachedStreamStatusFeedSubscription) {
	dfSub := &CachedStreamStatusFeedSubscription{}
	dfSub.ID = uuid.New().String()
	dfSub.Channel = make(chan *CachedStreamStatus)

	lc.streamStatusFeedSource._addSubscription(dfSub)

	return dfSub
}

func (lc *LiveCache) removeStreamStatusFeedSubscription(dfSub *CachedStreamStatusFeedSubscription) {
	lc.streamStatusFeedSource._removeSubscription(dfSub)
}

func (dfSource *t_streamStatusFeedSource) _addSubscription(dfSub *CachedStreamStatusFeedSubscription) {
	dfSource.subscriptionsMutex.Lock()
	defer dfSource.subscriptionsMutex.Unlock()

	dfSource.subscriptions = append(dfSource.subscriptions, dfSub)
	dfSource.subscriptionsIndex[dfSub.ID] = dfSub
}

func (dfSource *t_streamStatusFeedSource) _removeSubscription(dfSub *CachedStreamStatusFeedSubscription) {
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

		resizedSubs := make([]*CachedStreamStatusFeedSubscription, 0, subLen - 1 + stream_status_feed_subscription__extra_cap)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, dfSource.subscriptions[subIdx+1:subLen]...)

		dfSource.subscriptions = resizedSubs
		delete(dfSource.subscriptionsIndex, dfSub.ID)
	}

	// instead of dynamic dynamic buffer management
	// try high and low watermark indicator
}
