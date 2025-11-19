package util

import (
	"encoding/json"
)

const (
	BasicMessageType__Event  string = "_bev"
	BasicMessageType__ReqRep string = "_brr"
)

type BasicMessage struct {
	Type    string `json:"_bt"`
	Payload any    `json:"_bp"`
}

// type BasicEvent struct {
// 	Timestamp time.Time `json:"ts"`
// 	EventCode string    `json:"ev_code"`
// }

type BasicReqRep struct {
	RequestID   string `json:"_bid"`
	RequestCode string `json:"_brc"`
}

func NewMessage(msgType string, payload any) (BasicMessage) {
	return BasicMessage{msgType, payload}
}
func NewEventMessage(payload any) (BasicMessage) {
	return BasicMessage{BasicMessageType__Event, payload}
}
func NewReqRepMessage(payload any) (BasicMessage) {
	return BasicMessage{BasicMessageType__ReqRep, payload}
}

func NewMessageJson(msgType string, payload any) ([]byte) {
	msgBytes, err := json.Marshal(BasicMessage{msgType, payload})
	if err != nil {
		return []byte("_bmerr")
	}
	return msgBytes
}
func NewEventJson(payload any) ([]byte) {
	msgBytes, err := json.Marshal(BasicMessage{BasicMessageType__Event, payload})
	if err != nil {
		return []byte("_bmerr")
	}
	return msgBytes
}
func NewReqRepJson(payload any) ([]byte) {
	msgBytes, err := json.Marshal(BasicMessage{BasicMessageType__ReqRep, payload})
	if err != nil {
		return []byte("_bmerr")
	}
	return msgBytes
}

func ExtractBasicReqRep(reqBytes []byte) (BasicReqRep) {
	var reqStruct *BasicReqRep

	unmarshalErr := json.Unmarshal(reqBytes, &reqStruct)
	if unmarshalErr != nil {
		reqStruct = &BasicReqRep{"_b0", "_bumerr"}
	}

	return *reqStruct
}
