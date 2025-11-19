package pubsub

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"

	dilemmaKvmsg "noname001/dilemma/comm/zmqdep/kvmsg"

	deviceTyping "noname001/app/base/typing/device"
)

// ========================== VVV SUBSCRIPTION VVV ========================== //
type DeviceEventSubscription struct {
	id string
	MessageChannel chan *deviceTyping.LiveDeviceEvent
}

func (sub *DeviceSubscriber) SubscribeToDeviceEvent() (*DeviceEventSubscription) {
	subscription := &DeviceEventSubscription{
		id: uuid.New().String(),
		MessageChannel: make(chan *deviceTyping.LiveDeviceEvent),
	}

	sub.deviceEventSubscriptions[subscription.id] = subscription

	return subscription
}
func (sub *DeviceSubscriber) UnsubscribeFromDeviceEvent(subscription *DeviceEventSubscription) {
	_subscription, ok := sub.deviceEventSubscriptions[subscription.id]
	if ok {
		// close(_subscription.MessageChannel)
		delete(sub.deviceEventSubscriptions, _subscription.id)
	}
}
// ========================== ^^^ SUBSCRIPTION ^^^ ========================== //

// =========================== VVV SUBSCRIBER VVV =========================== //
func (sub *DeviceSubscriber) receiveDeviceEvent(msgJson string) {
	var (
		msg *deviceTyping.LiveDeviceEvent
		err error
	)

	err = json.Unmarshal([]byte(msgJson), &msg)
	if err != nil {
		// TODO
		sub.logger.Errorf("%s: receiveDeviceEvent - unmarshall err %s.", sub.logPrefix, err.Error())
		return
	}

	for _, subscription := range sub.deviceEventSubscriptions {
		subscription.MessageChannel <- msg
	}
}
// =========================== ^^^ SUBSCRIBER ^^^ =========================== //

// =========================== VVV PUBLISHER VVV ============================ //
func (pub *DevicePublisher) PublishDeviceEvent(evCode deviceTyping.LiveDeviceEventCode, nodeID, deviceID, deviceCode string, ttl int) (error) {
	var err error

	msgStruct := deviceTyping.LiveDeviceEvent{
		Timestamp: time.Now(),
		EventCode: evCode,

		NodeID    : nodeID,
		DeviceID  : deviceID,
		DeviceCode: deviceCode,
	}

	msgJson, err := json.Marshal(msgStruct)
	if err != nil {
		// TODO:
		pub.logger.Errorf("%s: PublishDeviceEvent - marshal err: %s.", pub.logPrefix, err.Error())
		return err
	}

	payload := HEADER__EVENT + HEADER_PAYLOAD_DELIM + string(msgJson)

	kvmsg := dilemmaKvmsg.NewKvmsg(0)
	kvmsg.SetKey(SUBTREE__DEVICE_EVENT)
	kvmsg.SetProp("ttl", strconv.Itoa(ttl))
	kvmsg.SetBody(payload)

	err = pub.pubClient.Publish(kvmsg)
	if err != nil {
		// TODO
		pub.logger.Errorf("%s: PublishDeviceEvent - publish err: %s.", pub.logPrefix, err.Error())
		return err
	}

	return nil
}
// =========================== ^^^ PUBLISHER ^^^ ============================ //
