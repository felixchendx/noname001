package mdpservice

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	SERVICE_CODE = "STRMPVDR"
	PROVIDER_ID_TEMPLATE = "%s::" + SERVICE_CODE

	DEF = "0.0.1"
	REQPARAM_DELIM = "|><|"

	DEFAULT_RETRY_BACKOFF = 10 * time.Second

	CMD_PING = string(iota + 1)
	CMD_SERVICE_INFO
	CMD_STREAM_SNAPSHOT_LIST
	CMD_STREAM_SNAPSHOT

	// CMD_TIME_SYNC
)

var (
	// TODO: def compat rules

	// TODO: register errors here, do not use fmt.Errorf() ?

	COMMANDS = map[string]string{
		CMD_PING                : "PING",
		CMD_SERVICE_INFO        : "SERVICE_INFO",
		CMD_STREAM_SNAPSHOT_LIST: "STREAM_SNAPSHOT_LIST",
		CMD_STREAM_SNAPSHOT     : "STREAM_SNAPSHOT",
	}
)

func StreamProviderID(nodeID string) (providerID string) {
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
