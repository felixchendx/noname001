package pubsub

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"

	dilemmaKvmsg "noname001/dilemma/comm/zmqdep/kvmsg"
)

// ========================== VVV SUBSCRIPTION VVV ========================== //
type NodeLivenessSubscription struct {
	id string
	MessageChannel chan *NodeLiveness
}

func (sub *NodeSubscriber) SubscribeToNodeLiveness() (*NodeLivenessSubscription) {
	subscription := &NodeLivenessSubscription{
		id: uuid.New().String(),
		MessageChannel: make(chan *NodeLiveness),
	}

	sub.nodeLivenessSubscriptions[subscription.id] = subscription

	return subscription
}
func (sub *NodeSubscriber) UnsubscribeFromNodeLiveness(subscription *NodeLivenessSubscription) {
	_subscription, ok := sub.nodeLivenessSubscriptions[subscription.id]
	if ok {
		// close(_subscription.MessageChannel)
		delete(sub.nodeLivenessSubscriptions, _subscription.id)
	}
}
// ========================== ^^^ SUBSCRIPTION ^^^ ========================== //

// =========================== VVV SUBSCRIBER VVV =========================== //
func (sub *NodeSubscriber) receiveNodeLiveness(msgJson string) {
	var (
		msg *NodeLiveness
		err error
	)

	err = json.Unmarshal([]byte(msgJson), &msg)
	if err != nil {
		// TODO
		sub.logger.Errorf("%s: receiveNodeLiveness - unmarshall err %s.", sub.logPrefix, err.Error())
		return
	}

	for _, subscription := range sub.nodeLivenessSubscriptions {
		subscription.MessageChannel <- msg
	}
}
// =========================== ^^^ SUBSCRIBER ^^^ =========================== //

// =========================== VVV PUBLISHER VVV ============================ //
type NodeLiveness struct {
	NodeID string
	LocalTime time.Time
}

func (pub *NodePublisher) PublishNodeLiveness(nodeID string, ttl int) (error) {
	var err error
	
	msgStruct := NodeLiveness{
		NodeID: nodeID,
		LocalTime: time.Now(),
	}

	msgJson, err := json.Marshal(msgStruct)
	if err != nil {
		// TODO:
		pub.logger.Errorf("%s: PublishNodeLiveness - marshal err: %s.", pub.logPrefix, err.Error())
		return err
	}

	payload := HEADER__LIVENESS + HEADER_PAYLOAD_DELIM + string(msgJson)
	
	kvmsg := dilemmaKvmsg.NewKvmsg(0)
	kvmsg.SetKey(SUBTREE__NODE_LIVENESS)
	kvmsg.SetProp("ttl", strconv.Itoa(ttl))
	kvmsg.SetBody(payload)

	err = pub.pubClient.Publish(kvmsg)
	if err != nil {
		// TODO
		pub.logger.Errorf("%s: PublishNodeLiveness - publish err: %s.", pub.logPrefix, err.Error())
		return err
	}

	return nil
}
// =========================== ^^^ PUBLISHER ^^^ ============================ //
