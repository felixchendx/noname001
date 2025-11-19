package mdpservice

import (
	"noname001/app/base/messaging"
	baseTyping "noname001/app/base/typing"
)

type DataHandlerIntface interface {
	ProvideServiceInfo() (*ServiceInfo, *messaging.Messages)

	ProvideDeviceSnapshotList()   ([]*baseTyping.BaseDeviceSnapshot, *messaging.Messages)
	ProvideDeviceSnapshot(string) (*baseTyping.BaseDeviceSnapshot, *messaging.Messages)
}
