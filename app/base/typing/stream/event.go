package stream

import (
	"time"
)

type LiveStreamEventCode string

const (
	LIVE_STREAM_EVENT_CODE__DEACTIVATED  LiveStreamEventCode = "ls:deactivated"
	LIVE_STREAM_EVENT_CODE__INIT_BEGIN   LiveStreamEventCode = "ls:init:begin"
	LIVE_STREAM_EVENT_CODE__INIT_FAIL    LiveStreamEventCode = "ls:init:fail"
	LIVE_STREAM_EVENT_CODE__INIT_OK      LiveStreamEventCode = "ls:init:ok"
	LIVE_STREAM_EVENT_CODE__RELOAD_BEGIN LiveStreamEventCode = "ls:reload:begin"
	LIVE_STREAM_EVENT_CODE__RELOAD_FAIL  LiveStreamEventCode = "ls:reload:fail"
	LIVE_STREAM_EVENT_CODE__RELOAD_OK    LiveStreamEventCode = "ls:reload:ok"
	LIVE_STREAM_EVENT_CODE__DESTROYED    LiveStreamEventCode = "ls:destroyed"
	LIVE_STREAM_EVENT_CODE__BG_FAIL      LiveStreamEventCode = "ls:bg:fail"
)

type LiveStreamEvent struct {
	Timestamp time.Time
	// EventID   string
	EventCode LiveStreamEventCode

	NodeID     string
	StreamID   string
	StreamCode string
}
