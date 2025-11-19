package util

import (
	"encoding/json"
)

const (
	BASIC_MESSAGE_TYPE__REQ_REP     string = "_brr2"
	BASIC_MESSAGE_TYPE__REQ_REP_ERR string = "_brr2err" // temp
)

type BasicReqHeader struct {
	Type    string `json:"_bt"`

	ReqID   string `json:"_bid"`
	ReqCode string `json:"_brc"`
}

type BasicReqMessage struct {
	Header  *BasicReqHeader `json:"_bh"`
	// Payload any            `json:"_bp"` // handle payload where the type def is
}

type BasicRepHeader struct {
	Type    string `json:"_bt"`

	ReqID   string `json:"_bid"`
	ReqCode string `json:"_brc"`

	// rep
}

type BasicRepMessage struct {
	Header  BasicRepHeader `json:"_bh"`
	Payload any            `json:"_bp"`
}

func ExtractBasicReqHeader(msgBytes []byte) (*BasicReqHeader) {
	var reqMessage *BasicReqMessage

	unmarshalErr := json.Unmarshal(msgBytes, &reqMessage)
	if unmarshalErr != nil {
		// TODO: logging
		reqMessage = &BasicReqMessage{
			&BasicReqHeader{BASIC_MESSAGE_TYPE__REQ_REP_ERR, "_b0", "_bumerr"},
		}
	}

	return reqMessage.Header
}

func DecodeTypedReqMessage(msgBytes []byte, typedReqMessage any) {
	unmarshalErr := json.Unmarshal(msgBytes, typedReqMessage)
	if unmarshalErr != nil {
		// TODO: logging
		typedReqMessage = &BasicReqMessage{
			&BasicReqHeader{BASIC_MESSAGE_TYPE__REQ_REP_ERR, "_b0", "_bumerr"},
		}
	}
}

func EncodeBasicRepMessage(reqHeader *BasicReqHeader, payload any) ([]byte) {
	repHeader := BasicRepHeader{BASIC_MESSAGE_TYPE__REQ_REP, reqHeader.ReqID, reqHeader.ReqCode}

	msg := &BasicRepMessage{repHeader, payload}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		// TODO: logging
		msgBytes, _ = json.Marshal(&BasicRepMessage{
			BasicRepHeader{BASIC_MESSAGE_TYPE__REQ_REP_ERR, "_b0", "_bmerr"},
			nil,
		})
	}

	return msgBytes
}
