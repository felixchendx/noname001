package hikvision

import (
	baseTyping "noname001/app/base/typing"

	appConstant "noname001/app/constant"
)

func (dev *HikvisionDevice) Init() {
	switch dev.cache.lDat.State {
	case baseTyping.DEVICE_LIVE_STATE__NEW:
		// pass

	default:
		return // can only be initialized once
	}

	switch dev.cache.pDat.State {
	case appConstant.ENTITY__STATE_ACTIVE:
		dev.doInit()
		dev.cron.Start()

	case appConstant.ENTITY__STATE_INACTIVE:
		fallthrough
	default:
		dev.doDeactivate()
		dev.cron.Stop()
	}
}

func (dev *HikvisionDevice) Reload() {
	switch dev.cache.lDat.State {
	case baseTyping.DEVICE_LIVE_STATE__INIT_BEGIN  :
		fallthrough // init is still running
	case baseTyping.DEVICE_LIVE_STATE__RELOAD_BEGIN:
		fallthrough // prev reload is still running
	case baseTyping.DEVICE_LIVE_STATE__DESTROY     :
		return      // reload what ?

	default:
		// pass
	}

	switch dev.cache.pDat.State {
	case appConstant.ENTITY__STATE_ACTIVE:
		dev.doReload()
		dev.cron.Start() // in case of state change

	case appConstant.ENTITY__STATE_INACTIVE:
		fallthrough
	default:
		dev.doDeactivate()
		dev.cron.Stop() // in case of state change
	}
}

func (dev *HikvisionDevice) Destroy() {
	switch dev.cache.lDat.State {
	case baseTyping.DEVICE_LIVE_STATE__DESTROY:
		return // can only be destroyed once

	default:
		// pass
	}

	dev.doDestroy()
	dev.cron.Stop()
}

func (dev *HikvisionDevice) PersistenceData() (baseTyping.BaseDevicePersistenceData) {
	return dev.cache.pDat
}

func (dev *HikvisionDevice) DeviceSnapshot() (*baseTyping.BaseDeviceSnapshot) {
	return &baseTyping.BaseDeviceSnapshot{
		dev.cache.pDat,
		dev.cache.lDat,
		dev.cache.opCap,
		dev.cache.hwDat,
	}
}

func (dev *HikvisionDevice) TempErrorDetails() (map[string]string) {
	return dev.cache.tempErrDetails
}
