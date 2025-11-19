package util

import (
	"encoding/json"
)

const (
	BASIC_MESSAGE_TYPE__DATA_FEED     string = "_bdf"
	BASIC_MESSAGE_TYPE__DATA_FEED_ERR string = "_bdferr" // temp
)

type BasicDataFeedHeader struct {
	Type  string `json:"_bt"`

	From  string `json:"_bfrom"`
	To    string `json:"_bto"`
	Topic string `json:"_btopic"`
}

type BasicDataFeedMessage struct {
	Header  BasicDataFeedHeader `json:"_bh"`
	Payload any                 `json:"_bp"`
}

func EncodeBasicDatafeedMessage(header BasicDataFeedHeader, payload any) ([]byte) {
	header.Type = BASIC_MESSAGE_TYPE__DATA_FEED
	msg := &BasicDataFeedMessage{
		header,
		payload,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		// TODO: logging
		header.Type = BASIC_MESSAGE_TYPE__DATA_FEED_ERR
		msgBytes, _ = json.Marshal(&BasicDataFeedMessage{
			header,
			nil,
		})
	}

	return msgBytes
}
