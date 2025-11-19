package typing

import (
	"time"
)

type NodeEventCode string

const (
	NODE_EVENT_CODE__READY     NodeEventCode = "n:ready"
	NODE_EVENT_CODE__IP_CHANGE NodeEventCode = "n:ip_change"
	NODE_EVENT_CODE__SHUTDOWN  NodeEventCode = "n:shutdown"

	NODE_EVENT_CODE__HEARTBEAT NodeEventCode = "n:heartbeat"
)

type NodeEvent struct {
	Timestamp time.Time
	// EventID   string
	EventCode NodeEventCode

	NodeID string
}
