package service

import (
	"noname001/app/base/messaging"

	baseTyping "noname001/app/base/typing"

	deviceMdp "noname001/app/node/comm/device/mdpservice"
)

// ============ VVV conform to deviceMdp.DataHandlerIntface VVV ============= //
func (svc *Service) ProvideServiceInfo() (*deviceMdp.ServiceInfo, *messaging.Messages) {
	messages := messaging.NewMessages()
	serviceInfo := &deviceMdp.ServiceInfo{
		Placeholder: "d-nwoo",
	}

	return serviceInfo, messages
}

func (svc *Service) ProvideDeviceSnapshotList() ([]*baseTyping.BaseDeviceSnapshot, *messaging.Messages) {
	messages := messaging.NewMessages()
	deviceSnapshotList := make([]*baseTyping.BaseDeviceSnapshot, 0)

	rawDeviceSnapshots := svc.coordinator.GetDeviceSnapshots()
	for _, deviceSnapshot := range rawDeviceSnapshots {
		deviceSnapshotList = append(deviceSnapshotList, svc.prepDeviceSnapshotForRemote(deviceSnapshot))
	}

	return deviceSnapshotList, messages
}

func (svc *Service) ProvideDeviceSnapshot(deviceCode string) (*baseTyping.BaseDeviceSnapshot, *messaging.Messages) {
	messages := messaging.NewMessages()

	deviceSnapshot := svc.coordinator.GetDeviceSnapshotByCode(deviceCode)
	if deviceSnapshot == nil {
		messages.AddError(CDT_ERR_90503.NewMessage(deviceCode))
		return nil, messages
	}

	return svc.prepDeviceSnapshotForRemote(deviceSnapshot), messages
}
// ============ ^^^ conform to deviceMdp.DataHandlerIntface ^^^ ============= //

func (svc *Service) prepDeviceSnapshotForRemote(_ds *baseTyping.BaseDeviceSnapshot) (*baseTyping.BaseDeviceSnapshot) {
	// TODO: time sensitive stuffs convert here
	// _deviceSnapshot.Live.LastSeen = _deviceSnapshot.Live.LastSeen

	return _ds
}
