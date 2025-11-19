package live

import (
	"context"
	"time"

	nodeTyping   "noname001/node/base/typing"
	deviceTyping "noname001/app/base/typing/device"
	streamTyping "noname001/app/base/typing/stream"
)

const (
	// assuming n inbound events per ticker duration
	// includes all events from 1 node (i.e. devices, streams)
	// node_action_buffer_cap int = 32

	// node_exec_tick_dur = 1000 * time.Millisecond

	node__external_activity_logs_limit = 20
	node__internal_activity_logs_limit = 10

	node__last_external_activity_threshold int64 = 60 * 3 // 3 mins

	node__ttl int64 = 60 * 60 * 24 * 1 // 1 day
)

func (lc *LiveCache) addNewNode(nodeSnapshot *nodeTyping.BaseNodeSnapshot) (*t_node) {
	node := &t_node{}
	node.id = nodeSnapshot.ID

	node.context, node.cancel = context.WithCancel(lc.context)

	node.nodeSnapshot = nodeSnapshot
	node.nodeStatus = nil
	node.nodeResource = nil

	node.mediaServer = lc.newMediaServer()
	node.deviceService = lc.newDeviceService()
	node.streamService = lc.newStreamService()

	node.nodeEventForwardingChan   = make(chan *nodeTyping.NodeEvent)
	node.deviceEventForwardingChan = make(chan *deviceTyping.LiveDeviceEvent)
	node.streamEventForwardingChan = make(chan *streamTyping.LiveStreamEvent)

	node.lastInternalActivityAt = time.Now()
	node.internalActivityLogs = make([]*t_internalActivityLog, 0, node__internal_activity_logs_limit + 1)

	node.lastExternalActivityAt = time.Now()
	node.externalActivityLogs = make([]*t_externalActivityLog, 0, node__external_activity_logs_limit + 1)

	node.staleSince = zero_time
	node.ttl = node__ttl

	lc.registerNode(node)
	go lc.nodeWorker(node)

	return node
}

func (lc *LiveCache) updateNodeData(node *t_node, nodeSnapshot *nodeTyping.BaseNodeSnapshot) {
	node.nodeSnapshot = nodeSnapshot

	node.defunctAt = zero_time
	node.defunctReason = ""

	node.logExternalActivity("data_update", nil)
}

func (lc *LiveCache) markNodeAsDefunct(node *t_node, reason string) {
	node.defunctAt = time.Now()
	node.defunctReason = reason

	node.logExternalActivity("defunct", []string{reason})
}

func (lc *LiveCache) updateNodeResource(node *t_node, nodeResource *nodeTyping.TempNodeSystemResourceSummary) {
	node.nodeResource = nodeResource

	// TODO:
}

func (lc *LiveCache) checkNodeStaleness(node *t_node, checkTime time.Time) {
	if (checkTime.Unix() - node.lastExternalActivityAt.Unix()) > node__last_external_activity_threshold {
		if node.staleSince.IsZero() {
			node.staleSince = checkTime
		}
		node.logInternalActivity("stale_check", "stale")

	} else {
		node.staleSince = zero_time
		node.logInternalActivity("stale_check", "ok")
	}
}

func (lc *LiveCache) isNodeExpired(node *t_node, checkTime time.Time) (bool) {
	var expired bool = false

	if node.staleSince.IsZero() {
		expired = false

	} else {
		expired = (checkTime.Unix() - node.lastExternalActivityAt.Unix()) > node.ttl
	}

	return expired
}

func (lc *LiveCache) removeNode(node *t_node) {
	node.cancel()
	lc.deregisterNode(node)
}

// -----------------------------------------------------------------------------
type t_node struct {
	id string

	// === data holder ===
	nodeSnapshot *nodeTyping.BaseNodeSnapshot
	nodeStatus   *t_nodeStatusInterpretation
	nodeResource *nodeTyping.TempNodeSystemResourceSummary

	mediaServer   *t_mediasrv
	deviceService *t_deviceService
	streamService *t_streamService
	// === data holder ===

	// === internal ===
	context context.Context
	cancel  context.CancelFunc

	nodeEventForwardingChan   chan *nodeTyping.NodeEvent
	deviceEventForwardingChan chan *deviceTyping.LiveDeviceEvent
	streamEventForwardingChan chan *streamTyping.LiveStreamEvent

	lastInternalActivityAt time.Time
	internalActivityLogs []*t_internalActivityLog

	lastExternalActivityAt time.Time
	externalActivityLogs []*t_externalActivityLog

	// TODO: aggregated activity log

	defunctAt     time.Time
	defunctReason string

	staleSince time.Time
	ttl int64
}

type t_nodeStatusInterpretation struct {
	timestamp time.Time

	textualIndicator  string
	visualIndicator   string
	auditoryIndicator string
}

func (node *t_node) logInternalActivity(activity, result string) {
	var naw = time.Now()

	node.lastInternalActivityAt = naw
	node.internalActivityLogs = append(node.internalActivityLogs, &t_internalActivityLog{
		ts: naw,
		activity: activity,
		result: result,
	})

	if len(node.internalActivityLogs) > node__internal_activity_logs_limit {
		node.internalActivityLogs = node.internalActivityLogs[1:len(node.internalActivityLogs)]
	}
}

func (node *t_node) logExternalActivity(activity string, extra []string) {
	var naw = time.Now()

	node.lastExternalActivityAt = naw
	node.externalActivityLogs = append(node.externalActivityLogs, &t_externalActivityLog{
		ts: naw,
		activity: activity,
		extra: extra,
	})

	if len(node.externalActivityLogs) > node__external_activity_logs_limit {
		node.externalActivityLogs = node.externalActivityLogs[1:len(node.externalActivityLogs)]
	}
}
