package mdpservice

import (
	"noname001/app/base/messaging"

	streamTyping "noname001/app/base/typing/stream"
)

type DataHandlerIntface interface {
	ProvideServiceInfo() (*ServiceInfo, *messaging.Messages)

	ProvideStreamSnapshotList() ([]*streamTyping.StreamSnapshot, *messaging.Messages)
	ProvideStreamSnapshot(streamCode string) (*streamTyping.StreamSnapshot, *messaging.Messages)
}
