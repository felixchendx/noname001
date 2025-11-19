package coordinator

import (
	"noname001/app/base/messaging"
	baseTyping  "noname001/app/base/typing"
	appConstant "noname001/app/constant"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

func (coord *Coordinator) GetDeviceSnapshots() ([]*baseTyping.BaseDeviceSnapshot) {
	deviceSnapshots := make([]*baseTyping.BaseDeviceSnapshot, 0)

	for _, liveDevice := range coord.liveDevicesOrdered {
		deviceSnapshots = append(deviceSnapshots, liveDevice.DeviceSnapshot())
	}

	return deviceSnapshots
}

func (coord *Coordinator) GetDeviceSnapshot(deviceID string) (*baseTyping.BaseDeviceSnapshot) {
	liveDevice, found := coord.liveDevices[deviceID]
	if !found {
		return nil
	}

	return liveDevice.DeviceSnapshot()
}

func (coord *Coordinator) GetDeviceSnapshotByCode(deviceCode string) (*baseTyping.BaseDeviceSnapshot) {
	liveDeviceID, idFound := coord.liveDevicesCode[deviceCode]
	if !idFound {
		return nil
	}

	liveDevice, found := coord.liveDevices[liveDeviceID]
	if !found {
		return nil
	}

	return liveDevice.DeviceSnapshot()
}

func (coord *Coordinator) GetTempErrorDetails(deviceID string) (map[string]string) {
	liveDevice, found := coord.liveDevices[deviceID]
	if !found {
		return make(map[string]string)
	}

	return liveDevice.TempErrorDetails()
}

func (coord *Coordinator) TempReloadDevice(deviceID string) (ok bool) {
	liveDevice, found := coord.liveDevices[deviceID]
	if !found {
		return false
	}

	liveDevice.Reload()

	return ok
}


func (coord *Coordinator) GetStreamInfo(deviceID, channelID string, streamType appConstant.BrandStreamType) (*liveBase.StreamInfo, *messaging.Messages) {
	messages := messaging.NewMessages()

	liveDevice, ok := coord.liveDevices[deviceID]
	if !ok {
		messages.AddError(SVC_ERR_90501.NewMessage())
		return nil, messages
	}

	streamInfo, deviceResponse := liveDevice.GetStreamInfo(channelID, streamType)
	if deviceResponse.IsConsideredError() {
		switch {
		case deviceResponse.IsGoError():
			messages.AddError(SVC_ERR_90502.NewMessage(deviceResponse.GoError().Error()))
		case deviceResponse.IsAPIError():
			messages.AddError(SVC_ERR_90503.NewMessage(deviceResponse.APIError().SimpleError()))
		}
		return nil, messages
	}

	return streamInfo, messages
}
