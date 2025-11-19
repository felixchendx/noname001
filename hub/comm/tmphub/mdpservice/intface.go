package mdpservice

import (
	"noname001/app/base/messaging"
)

type DataHandlerIntface interface {
	ProvideNodeList() ([]*NodeInfo, *messaging.Messages)
}
