package service

import (
	"noname001/app/base/messaging"
	baseTyping  "noname001/app/base/typing"
	appConstant "noname001/app/constant"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

func (svc *Service) GetDeviceSnapshots() ([]*baseTyping.BaseDeviceSnapshot) {
	return svc.coordinator.GetDeviceSnapshots()
}

func (svc *Service) GetDeviceSnapshot(deviceID string) (*baseTyping.BaseDeviceSnapshot, *messaging.Messages) {
	messages := messaging.NewMessages()

	deviceSnapshot := svc.coordinator.GetDeviceSnapshot(deviceID)
	if deviceSnapshot == nil {
		messages.AddError(CDT_ERR_90502.NewMessage(deviceID))
		return nil, messages
	}

	return deviceSnapshot, messages
}

func (svc *Service) GetDeviceSnapshotByCode(deviceCode string) (*baseTyping.BaseDeviceSnapshot, *messaging.Messages) {
	messages := messaging.NewMessages()

	deviceSnapshot := svc.coordinator.GetDeviceSnapshotByCode(deviceCode)
	if deviceSnapshot == nil {
		messages.AddError(CDT_ERR_90503.NewMessage(deviceCode))
		return nil, messages
	}

	return deviceSnapshot, messages
}

func (svc *Service) GetTempErrorDetails(deviceID string) (map[string]string) {
	return svc.coordinator.GetTempErrorDetails(deviceID)
}

func (svc *Service) ReloadDevice(deviceID string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	ok := svc.coordinator.TempReloadDevice(deviceID)
	if !ok {
		messages.AddError(CDT_ERR_90502.NewMessage(deviceID))
		return messages
	}

	return messages
}

func (svc *Service) GetStreamInfo(deviceIder *DeviceIdentifier, channelID string, streamType appConstant.BrandStreamType) (*liveBase.StreamInfo, *messaging.Messages) {
	messages := messaging.NewMessages()

	deviceDE, messages01 := svc._getDeviceWithIdentifier(deviceIder)
	if messages01.HasError() {
		messages.Append(messages01)
		return nil, messages
	}

	if deviceDE.State == appConstant.ENTITY__STATE_INACTIVE {
		messages.AddError(CDT_ERR_14501.NewMessage(deviceDE.Code))
		return nil, messages
	}

	streamInfo, messages02 := svc.coordinator.GetStreamInfo(deviceDE.ID, channelID, streamType)
	if messages02.HasError() {
		messages.AddError(CDT_ERR_14502.NewMessage(deviceDE.Code))
		messages.Append(messages02)
		return nil, messages
	}

	return streamInfo, messages
}


func (svc *Service) _getDeviceWithIdentifier(deviceIder *DeviceIdentifier) (*DeviceDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	if deviceIder == nil || (deviceIder.ID == "" && deviceIder.Code == "") {
		messages.AddError(CDT_ERR_90501.NewMessage())
		return nil, messages
	}

	var deviceDE *DeviceDE
	switch {
	case deviceIder.ID != "":
		pe, dbev := svc.store.DB.Device__Get(deviceIder.ID)
		if dbev.IsError() {
			messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
			return nil, messages
		}
		if pe == nil {
			messages.AddError(CDT_ERR_90502.NewMessage(deviceIder.ID))
			return nil, messages
		}

		deviceDE = (&DeviceDE{}).fromPE(pe)

	case deviceIder.Code != "":
		pe, dbev := svc.store.DB.Device__GetByCode(deviceIder.Code)
		if dbev.IsError() {
			messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
			return nil, messages
		}
		if pe == nil {
			messages.AddError(CDT_ERR_90503.NewMessage(deviceIder.Code))
			return nil, messages
		}

		deviceDE = (&DeviceDE{}).fromPE(pe)
	}

	return deviceDE, messages
}
