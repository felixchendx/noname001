package event

import (
	"sync"

	"github.com/google/uuid"

	nodeTyping "noname001/node/base/typing"
)

const (
	node_event_subscription__extra_cap int = 1
)

type NodeEventCode = nodeTyping.NodeEventCode
type NodeEvent     = nodeTyping.NodeEvent

type NodeEventSubscription struct {
	id string

	Channel chan NodeEvent
}

type t_nodeEventSource struct {
	sourceChan chan NodeEvent

	// === subscriptionsMutex ===
	subscriptionsMutex sync.Mutex
	subscriptions      []*NodeEventSubscription
	// === subscriptionsMutex ===

	// closerChan chan *NodeEventSubscription
}

func (evHub *EventHub) newNodeEventSource() (*t_nodeEventSource) {
	evSource := &t_nodeEventSource{}
	evSource.sourceChan = make(chan NodeEvent)

	evSource.subscriptions = make([]*NodeEventSubscription, 0, node_event_subscription__extra_cap)

	go func() {
		evSourceLoop:
		for {
			select {
			case <- evHub.context.Done():
				break evSourceLoop

			case ev := <- evSource.sourceChan:
				for _, evSub := range evSource.subscriptions {
					evSub.Channel <- ev
				}
			}
		}
	}()

	return evSource
}

func (evHub *EventHub) PublishNodeEvent(evCode NodeEventCode, nodeID string) {
	evHub.nodeEventSource.sourceChan <- NodeEvent{
		EventCode: evCode,
		NodeID   : nodeID,
	}
}

func (evHub *EventHub) NewNodeEventSubscription() (*NodeEventSubscription) {
	evSub := &NodeEventSubscription{}
	evSub.id = uuid.New().String()
	evSub.Channel = make(chan NodeEvent)

	evHub.nodeEventSource._addSubscription(evSub)

	return evSub
}

func (evHub *EventHub) RemoveNodeEventSubscription(evSub *NodeEventSubscription) {
	evHub.nodeEventSource._removeSubscription(evSub)
}

func (evSource *t_nodeEventSource) _addSubscription(evSub *NodeEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	evSource.subscriptions = append(evSource.subscriptions, evSub)

	evSource.subscriptionsMutex.Unlock()
}

func (evSource *t_nodeEventSource) _removeSubscription(evSub *NodeEventSubscription) {
	evSource.subscriptionsMutex.Lock()

	subIdx := -1

	for _subIdx, _evSub := range evSource.subscriptions {
		if _evSub.id == evSub.id {
			subIdx = _subIdx
			break
		}
	}

	if subIdx != -1 {
		subLen := len(evSource.subscriptions)

		resizedSubs := make([]*NodeEventSubscription, 0, subLen - 1 + node_event_subscription__extra_cap)
		resizedSubs = append(resizedSubs, evSource.subscriptions[0:subIdx]...)
		resizedSubs = append(resizedSubs, evSource.subscriptions[subIdx+1:subLen]...)

		evSource.subscriptions = resizedSubs
	}

	evSource.subscriptionsMutex.Unlock()
}
