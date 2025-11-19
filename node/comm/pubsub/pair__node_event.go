package pubsub

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"

	dilemmaKvmsg "noname001/dilemma/comm/zmqdep/kvmsg"

	nodeTyping "noname001/node/base/typing"
)

// ========================== VVV SUBSCRIPTION VVV ========================== //
type NodeEventSubscription struct {
	id string
	MessageChannel chan *nodeTyping.NodeEvent
}

func (sub *NodeSubscriber) SubscribeToNodeEvent() (*NodeEventSubscription) {
	subscription := &NodeEventSubscription{
		id: uuid.New().String(),
		MessageChannel: make(chan *nodeTyping.NodeEvent),
	}

	sub.nodeEventSubscriptions[subscription.id] = subscription

	return subscription
}
func (sub *NodeSubscriber) UnsubscribeFromNodeEvent(subscription *NodeEventSubscription) {
	_subscription, ok := sub.nodeEventSubscriptions[subscription.id]
	if ok {
		// close(_subscription.MessageChannel)
		delete(sub.nodeEventSubscriptions, _subscription.id)
	}
}
// ========================== ^^^ SUBSCRIPTION ^^^ ========================== //

// =========================== VVV SUBSCRIBER VVV =========================== //
func (sub *NodeSubscriber) receiveNodeEvent(msgJson string) {
	var (
		msg *nodeTyping.NodeEvent
		err error
	)

	err = json.Unmarshal([]byte(msgJson), &msg)
	if err != nil {
		// TODO
		sub.logger.Errorf("%s: receiveNodeEvent - unmarshall err %s.", sub.logPrefix, err.Error())
		return
	}

	for _, subscription := range sub.nodeEventSubscriptions {
		subscription.MessageChannel <- msg
	}
}
// =========================== ^^^ SUBSCRIBER ^^^ =========================== //

// =========================== VVV PUBLISHER VVV ============================ //
func (pub *NodePublisher) PublishNodeEvent(evCode nodeTyping.NodeEventCode, nodeID string, ttl int) (error) {
	var err error

	msgStruct := nodeTyping.NodeEvent{
		Timestamp: time.Now(),
		EventCode: evCode,

		NodeID: nodeID,
	}

	msgJson, err := json.Marshal(msgStruct)
	if err != nil {
		// TODO:
		pub.logger.Errorf("%s: PublishNodeEvent - marshal err: %s.", pub.logPrefix, err.Error())
		return err
	}

	payload := HEADER__EVENT + HEADER_PAYLOAD_DELIM + string(msgJson)
	
	kvmsg := dilemmaKvmsg.NewKvmsg(0)
	kvmsg.SetKey(SUBTREE__NODE_EVENT)
	kvmsg.SetProp("ttl", strconv.Itoa(ttl))
	kvmsg.SetBody(payload)

	err = pub.pubClient.Publish(kvmsg)
	if err != nil {
		// TODO
		pub.logger.Errorf("%s: PublishNodeEvent - publish err: %s.", pub.logPrefix, err.Error())
		return err
	}

	return nil
}
// =========================== ^^^ PUBLISHER ^^^ ============================ //
