package intface

import (
	"time"

	streamTyping "noname001/app/base/typing/stream"
)

type CachedStreamFilter struct {}
type CachedStream struct {
	ID   string
	Code string

	StreamSnapshot streamTyping.StreamSnapshot
	StreamStatus   *CachedStreamStatus

	RelayPathName string
	SourceStream  string

	LastActivityAt time.Time
}

type CachedStreamStatus struct {
	StreamID   string

	NodeID     string
	StreamCode string
	Timestamp  time.Time

	TextualIndicator  string
	VisualIndicator   string
	AuditoryIndicator string
}

type CachedStreamEvent struct {
	OriginalStreamEvent *streamTyping.LiveStreamEvent
}


type CachedStreamStatusFeedSubscription struct {
	ID string

	Channel chan *CachedStreamStatus
}

type CachedStreamEventSubscription struct {
	ID string

	Channel chan CachedStreamEvent
}
