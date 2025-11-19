package mdpservice

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	SERVICE_CODE = "NODEPVDR"
	PROVIDER_ID_TEMPLATE = "%s::" + SERVICE_CODE

	DEF = "0.0.1"
	REQPARAM_DELIM = "|><|"

	DEFAULT_RETRY_BACKOFF = 10 * time.Second

	CMD__PING = string(iota + 1)
	CMD__NODE_SNAPSHOT
	CMD__NODE_RESOURCE

	// CMD_TIME_SYNC
)

var (
	// TODO: def compat rules

	// TODO: register errors here, do not use fmt.Errorf() ?

	COMMANDS = map[string]string{
		CMD__PING         : "PING",
		CMD__NODE_SNAPSHOT: "SNAPSHOT",
		CMD__NODE_RESOURCE: "RESOURCE",
	}
)

func nodeProviderID(nodeID string) (providerID string) {
	return fmt.Sprintf(PROVIDER_ID_TEMPLATE, nodeID)
}

type ErrorReply struct {
	Status  string
	Message string

	Data    any
}
func SerializedErrorReply(message string) (string) {
	repBytes, _ := json.Marshal(ErrorReply{"error", message, nil})
	return string(repBytes)
}
