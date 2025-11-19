package intface

import (
	"time"

	nodeTyping "noname001/node/base/typing"
)

type CachedNodeFilter struct {}
type CachedNode struct {
	LastActivityAt time.Time

	NodeSnapshot nodeTyping.BaseNodeSnapshot
	NodeStatus   *CachedNodeStatus
	NodeResource *CachedNodeResource

	// // TODO: services, in nodeSnapshot

	CachedDevices []*CachedDevice
	CachedStreams []*CachedStream

	DevicesXStreams map[string][]string // map[_deviceCode_][]_streamCode_
}
// make this optional
// func (nc NodeCache) GenerateDeviceXStreams() {}

type CachedNodeStatus struct {
	ID        string
	Timestamp time.Time

	TextualIndicator  string
	VisualIndicator   string
	AuditoryIndicator string
}

type CachedNodeResource struct {
	ID string

	NodeResource nodeTyping.TempNodeSystemResourceSummary
}

type CachedNodeEvent struct {
	OriginalNodeEvent *nodeTyping.NodeEvent
}

type CachedNodeStatusFeedSubscription struct {
	ID string

	Channel chan *CachedNodeStatus
}

type CachedNodeResourceFeedSubscription struct {
	ID string

	Channel chan *CachedNodeResource
}

type CachedNodeEventSubscription struct {
	ID string

	Channel chan CachedNodeEvent
}
