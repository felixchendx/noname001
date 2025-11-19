package pubsub

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"

	dilemmaKvmsg "noname001/dilemma/comm/zmqdep/kvmsg"

	streamTyping "noname001/app/base/typing/stream"
)

// ========================== VVV SUBSCRIPTION VVV ========================== //
type StreamEventSubscription struct {
	id string
	MessageChannel chan *streamTyping.LiveStreamEvent
}

func (sub *StreamSubscriber) SubscribeToStreamEvent() (*StreamEventSubscription) {
	subscription := &StreamEventSubscription{
		id: uuid.New().String(),
		MessageChannel: make(chan *streamTyping.LiveStreamEvent),
	}

	sub.streamEventSubscriptions[subscription.id] = subscription

	return subscription
}
func (sub *StreamSubscriber) UnsubscribeFromStreamEvent(subscription *StreamEventSubscription) {
	_subscription, ok := sub.streamEventSubscriptions[subscription.id]
	if ok {
		// close(_subscription.MessageChannel)
		delete(sub.streamEventSubscriptions, _subscription.id)
	}
}
// ========================== ^^^ SUBSCRIPTION ^^^ ========================== //

// =========================== VVV SUBSCRIBER VVV =========================== //
func (sub *StreamSubscriber) receiveStreamEvent(msgJson string) {
	var (
		msg *streamTyping.LiveStreamEvent
		err error
	)

	err = json.Unmarshal([]byte(msgJson), &msg)
	if err != nil {
		// TODO
		sub.logger.Errorf("%s: receiveStreamEvent - unmarshall err %s.", sub.logPrefix, err.Error())
		return
	}

	for _, subscription := range sub.streamEventSubscriptions {
		subscription.MessageChannel <- msg
	}
}
// =========================== ^^^ SUBSCRIBER ^^^ =========================== //

// =========================== VVV PUBLISHER VVV ============================ //
func (pub *StreamPublisher) PublishStreamEvent(evCode streamTyping.LiveStreamEventCode, nodeID, streamID, streamCode string, ttl int) (error) {
	var err error

	msgStruct := streamTyping.LiveStreamEvent{
		Timestamp: time.Now(),
		EventCode: evCode,

		NodeID    : nodeID,
		StreamID  : streamID,
		StreamCode: streamCode,
	}

	msgJson, err := json.Marshal(msgStruct)
	if err != nil {
		// TODO:
		pub.logger.Errorf("%s: PublishStreamEvent - marshal err: %s.", pub.logPrefix, err.Error())
		return err
	}

	payload := HEADER__EVENT + HEADER_PAYLOAD_DELIM + string(msgJson)
	
	kvmsg := dilemmaKvmsg.NewKvmsg(0)
	kvmsg.SetKey(SUBTREE__STREAM_EVENT)
	kvmsg.SetProp("ttl", strconv.Itoa(ttl))
	kvmsg.SetBody(payload)

	err = pub.pubClient.Publish(kvmsg)
	if err != nil {
		// TODO
		pub.logger.Errorf("%s: PublishStreamEvent - publish err: %s.", pub.logPrefix, err.Error())
		return err
	}

	return nil
}
// =========================== ^^^ PUBLISHER ^^^ ============================ //
