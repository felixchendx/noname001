package base

import (
	baseTyping "noname001/app/base/typing"
	appConstant "noname001/app/constant"
)

type LiveDeviceIntface interface {
	Init()
	Reload()
	// PatchAndReload(*DevicePatchParams) // do not generalize until more brand is explored
	Destroy()

	// implement these data fragments when necessary
	PersistenceData() baseTyping.BaseDevicePersistenceData
	// LiveData()        (baseTyping.BaseDeviceLiveData)
	// OpCap()           (baseTyping.BaseDeviceOpCap)
	// HardwareData()    (baseTyping.BaseDeviceHardwareData)
	DeviceSnapshot() *baseTyping.BaseDeviceSnapshot

	TempErrorDetails() map[string]string

	GetStreamInfo(
		channelID string,
		streamType appConstant.BrandStreamType,
	) (
		streamInfo *StreamInfo,
		deviceResponse DeviceResponseIntface, // TODO: terminate internally
	)
}
