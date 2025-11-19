package route

import (
	"time"

	wsUtil "noname001/web/base/comm/ws/util"
)

const (
	wd__wsRequestCode__wallInfo     string = "/wall/info"
	wd__wsRequestCode__wallItemInfo string = "/wall/item/info"

	wd__wsRequestCode__nodeInfoListing   string = "/node/info/listing"
	wd__wsRequestCode__streamInfoListing string = "/stream/info/listing"
)

type wd__wsReply__wallInfo struct {
	wsUtil.BasicReqRep

	WallInfo *wd__wallInfo `json:"wall_info"`
}

type wd__wsRequest__wallItemInfo struct {
	wsUtil.BasicReqRep

	WallItemID string `json:"wall_item_id"`
}
type wd__wsReply__wallItemInfo struct {
	wsUtil.BasicReqRep

	WallItemInfo *wd__wallItemInfo `json:"wall_item_info"`
}

type wd__wallInfo struct {
	LayoutCode string `json:"layout_code"`

	Items []*wd__wallItemInfo `json:"items"`
}
type wd__wallItemInfo struct {
	ID    string `json:"id"`
	Index int    `json:"index"`

	SourceNode   string `json:"source_node"`
	SourceStream string `json:"source_stream"`

	NodeInfo   *wd__nodeInfo   `json:"node_info"`
	StreamInfo *wd__streamInfo `json:"stream_info"`
}


type wd__wsReply__nodeInfoListing struct {
	wsUtil.BasicReqRep

	NodeInfoListing []*wd__nodeInfo `json:"node_info_listing"`
}

type wd__nodeInfo struct {
	ID    string `json:"id"`
	State string `json:"state"`

	StreamCount    int       `json:"stream_count"`
	LastActivityAt time.Time `json:"last_activity_at"`
}


type wd__wsRequest__streamInfoListing struct {
	wsUtil.BasicReqRep

	NodeID string `json:"node_id"`
}
type wd__wsReply__streamInfoListing struct {
	wsUtil.BasicReqRep

	StreamInfoListing []*wd__streamInfo `json:"stream_info_listing"`
}

type wd__streamInfo struct {
	ID   string `json:"id"`
	Code string `json:"code"`

	SourceType string `json:"source_type"`

	StreamerState         string `json:"streamer_state"`
	EstimatedVideoBitrate int    `json:"estimated_video_bitrate"`

	LastActivityAt time.Time `json:"last_activity_at"`

	PreviewURL string `json:"preview_url"`
}
