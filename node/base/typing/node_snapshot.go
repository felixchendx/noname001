package typing

import (
	"time"
)

type NodeState string

const (
	NODE_STATE__INIT         NodeState = "n_s:init"
	NODE_STATE__START        NodeState = "n_s:start"
	NODE_STATE__READY        NodeState = "n_s:ready"
	NODE_STATE__ABORT        NodeState = "n_s:abort"
	NODE_STATE__STOP         NodeState = "n_s:stop"
	NODE_STATE__SHUTDOWN     NodeState = "n_s:shutdown"
	// TODO RECONNECTED
	NODE_STATE__DISCONNECTED NodeState = "n_s:disconnected"
)

type BaseNodeSnapshot struct {
	ID    string
	Name  string
	State NodeState

	IPs             []string
	LastIPHistoryTs time.Time

	LocalTime   time.Time
	Timezone    string

	HubSnapshot *BaseHubSnapshot
	AppSnapshot *BaseAppSnapshot
}

type BaseHubSnapshot struct {

}

type BaseAppSnapshot struct {
	ModuleStates map[string]string

	TempMediasrvSnapshot *TempMediasrvSnapshot
}

type TempMediasrvSnapshot struct {
	Ports     map[string]string
	AuthnPair string
}
