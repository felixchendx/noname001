package mdpservice

import (
	"encoding/json"
	"time"
)

const (
	PROVIDER_ID = "TMPHUB::TMPHUBPVDR"
    //             ^^^^^^ fixed marker for this provider

	DEF = "0.0.1"
	REQPARAM_DELIM = "|><|"

	DEFAULT_RETRY_BACKOFF = 10 * time.Second

	CMD_PING = string(iota + 1)
	CMD_NODE_LIST
)

var (
	// TODO: def compat rules

	// TODO: register errors here, do not use fmt.Errorf() ?

	COMMANDS = map[string]string{
		CMD_PING:      "PING",
		CMD_NODE_LIST: "NODE_LIST",
	}
)

type ErrorReply struct {
	Status  string
	Message string

	Data    any
}
func SerializedErrorReply(message string) (string) {
	repBytes, _ := json.Marshal(ErrorReply{"error", message, nil})
	return string(repBytes)
}
