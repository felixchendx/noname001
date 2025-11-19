package route

import (
	wsUtil "noname001/web/base/comm/ws/util"

	baseTyping "noname001/app/base/typing"
)

const (
	sgd__wsRequestCode__streamItemListing string = "/stream-item/listing"
	// sgd__wsRequestCode__streamItemSetActive string = "/stream-item/set-active"
	// sgd__wsRequestCode__streamItemSetInactive string = "/stream-item/set-inactive"

	sgd__wsRequestCode__deviceSnapshotListing string = "/device/snapshot/listing"
)

type sgd__wsReply__streamItemListing struct {
	wsUtil.BasicReqRep

	Items []sgd__streamItem `json:"items"`
}

type sgd__wsReply__deviceSnapshotListing struct {
	wsUtil.BasicReqRep

	DeviceSnapshots []*baseTyping.BaseDeviceSnapshot `json:"device_snapshots"`
}

type sgd__streamItem struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	State string `json:"state"`
	Note  string `json:"note"`

	SourceType       string `json:"source_type"`
	DeviceCode       string `json:"device_code"`
	DeviceChannelID  string `json:"device_channel_id"`
	DeviceStreamType string `json:"device_stream_type"`
	ExternalURL      string `json:"external_url"`
	Filepath         string `json:"filepath"`

	StreamURL string `json:"stream_url"`
}
