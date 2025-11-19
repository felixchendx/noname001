package mdpservice

import (
	nodeTyping "noname001/node/base/typing"

	"noname001/app/base/messaging"
)

type DataHandlerIntface interface {
	ProvideNodeSnapshot() (*nodeTyping.BaseNodeSnapshot, *messaging.Messages)
	ProvideNodeResource() (*nodeTyping.TempNodeSystemResourceSummary, *messaging.Messages)
}
