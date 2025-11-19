package coordinator

import (
	appConstant "noname001/app/constant"

	localEv "noname001/dilemma/event"

	deviceStore "noname001/app/module/common/device/store/sqlite"

	liveBase      "noname001/app/module/common/device/coordinator/live/base"
	liveDahua     "noname001/app/module/common/device/coordinator/live/dahua"
	liveHikvision "noname001/app/module/common/device/coordinator/live/hikvision"

	panasonicNetCam "noname001/app/module/common/device/coordinator/live/panasonic/netcam"
)

// TODO: revisit all messages when working on bg event stuffs

func (coord *Coordinator) initLiveDevice(deviceID string) {
	bgev := localEv.LocalEventHub().NewBgEvent("initLiveDevice")

	devicePE := coord._retrieveDeviceDataFromDB(bgev, deviceID)
	if bgev.Messages.HasError() {
		return
	}

	decryptedPassword := coord._decryptDevicePassword(bgev, devicePE.Password)
	if bgev.Messages.HasError() {
		return
	}

	liveDevice := coord._newLiveDevice(bgev, devicePE, decryptedPassword)
	if bgev.Messages.HasError() {
		return
	}

	coord.addLiveDeviceToRegistry(deviceID, liveDevice)
	liveDevice.Init()

	_ = bgev // silent succ
}

func (coord *Coordinator) reloadLiveDevice(deviceID string) {
	bgev := localEv.LocalEventHub().NewBgEvent("reloadLiveDevice")

	devicePE := coord._retrieveDeviceDataFromDB(bgev, deviceID)
	if bgev.Messages.HasError() {
		return
	}

	decryptedPassword := coord._decryptDevicePassword(bgev, devicePE.Password)
	if bgev.Messages.HasError() {
		return
	}

	existingLiveDevice, alreadyExist := coord.liveDevices[deviceID]
	if alreadyExist {
		prevBrand := existingLiveDevice.PersistenceData().Brand

		if prevBrand != devicePE.Brand {
			coord.destroyLiveDevice(devicePE.ID)
			coord.initLiveDevice(devicePE.ID)

		} else {
			coord._patchAndReloadLiveDevice(bgev, devicePE, decryptedPassword)
			if bgev.Messages.HasError() {
				return
			}
		}
	}

	_ = bgev // silent succ
}

func (coord *Coordinator) destroyLiveDevice(deviceID string) {
	bgev := localEv.LocalEventHub().NewBgEvent("destroyLiveDevice")

	existingLiveDevice, alreadyExist := coord.liveDevices[deviceID]
	if !alreadyExist {
		return
	}

	existingLiveDevice.Destroy()
	coord.removeLiveDeviceFromRegistry(deviceID)

	_ = bgev // succ
}

func (coord *Coordinator) _retrieveDeviceDataFromDB(bgev *localEv.BgEvent, deviceID string) (*deviceStore.DevicePE) {
	devicePE, dbev := coord.store.DB.Device__Get(deviceID)
	if dbev.IsError() {
		bgev.Messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		bgev.Messages.AddError(DVC_ERR_11001.NewMessage(deviceID))
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return nil
	}
	if devicePE == nil {
		bgev.Messages.AddError(DVC_ERR_11002.NewMessage(deviceID))
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return nil
	}

	return devicePE
}

func (coord *Coordinator) _decryptDevicePassword(bgev *localEv.BgEvent, encryptedPassword string) (string) {
	decryptedPassword, err := coord.secBundle.Decrypt(encryptedPassword, "device")
	if err != nil {
		bgev.Messages.AddError(DVC_WRN_80501.NewMessage())
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return ""
	}

	return decryptedPassword
}

func (coord *Coordinator) _newLiveDevice(
	bgev              *localEv.BgEvent,
	devicePE          *deviceStore.DevicePE,
	decryptedPassword string,
) (
	liveDevice liveBase.LiveDeviceIntface,
) {
	var err error

	switch devicePE.Brand {
	case string(appConstant.BRAND__DAHUA):
		liveDevice, err = liveDahua.NewDevice(&liveDahua.DahuaDeviceParams{
			coord.context,
			coord.logger, coord.logPrefix,
			coord.evHub,
			coord.commBundle,
			coord.timezone,

			devicePE.ID, devicePE.Code, devicePE.Name, devicePE.State, devicePE.Brand,
			devicePE.Protocol, devicePE.Hostname, devicePE.Port, devicePE.Username, decryptedPassword,
			devicePE.FallbackRTSPPort,
		})

	case string(appConstant.BRAND__HIKVISION):
		liveDevice, err = liveHikvision.NewDevice(&liveHikvision.HikvisionDeviceParams{
			coord.context,
			coord.logger, coord.logPrefix,
			coord.evHub,
			coord.commBundle,
			coord.timezone,

			devicePE.ID, devicePE.Code, devicePE.Name, devicePE.State, devicePE.Brand,
			devicePE.Protocol, devicePE.Hostname, devicePE.Port, devicePE.Username, decryptedPassword,
			devicePE.FallbackRTSPPort,
		})

	case string(appConstant.BRAND__PANASONIC_NETCAM):
		liveDevice, err = panasonicNetCam.NewPanasonicNetworkCamera(&panasonicNetCam.PanasonicNetworkCameraParams{
			coord.context,
			coord.logger, coord.logPrefix,
			coord.evHub,
			coord.commBundle,
			coord.timezone,

			devicePE.ID, devicePE.Code, devicePE.Name, devicePE.State, devicePE.Brand,
			devicePE.Protocol, devicePE.Hostname, devicePE.Port, devicePE.Username, decryptedPassword,
			devicePE.FallbackRTSPPort,
		})

	default:
		bgev.Messages.AddError(DVC_ERR_11299.NewMessage(devicePE.Code, devicePE.Brand))
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return nil
	}

	if err != nil {
		bgev.Messages.AddError(DVC_ERR_11201.NewMessage(devicePE.Code, err.Error()))
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return nil
	}

	return liveDevice
}

func (coord *Coordinator) _patchAndReloadLiveDevice(
	bgev              *localEv.BgEvent,
	devicePE          *deviceStore.DevicePE,
	decryptedPassword string,
) {
	var err error

	existingLiveDevice, alreadyExist := coord.liveDevices[devicePE.ID]
	if !alreadyExist {
		return
	}

	switch devicePE.Brand {
	case string(appConstant.BRAND__DAHUA):
		dahuaDevice, assertionOk := existingLiveDevice.(*liveDahua.DahuaDevice)
		if assertionOk {
			err = dahuaDevice.PatchAndReload(&liveDahua.DahuaDevicePatchParams{
				devicePE.Name, devicePE.State, devicePE.Brand,
				devicePE.Protocol, devicePE.Hostname, devicePE.Port, devicePE.Username, decryptedPassword,
				devicePE.FallbackRTSPPort,
			})
		} else {
			// warning ?
		}

	case string(appConstant.BRAND__HIKVISION):
		hikDevice, assertionOk := existingLiveDevice.(*liveHikvision.HikvisionDevice)
		if assertionOk {
			err = hikDevice.PatchAndReload(&liveHikvision.HikvisionDevicePatchParams{
				devicePE.Name, devicePE.State, devicePE.Brand,
				devicePE.Protocol, devicePE.Hostname, devicePE.Port, devicePE.Username, decryptedPassword,
				devicePE.FallbackRTSPPort,
			})
		} else {
			// warning ?
		}

	case string(appConstant.BRAND__PANASONIC_NETCAM):
		panNetCam, assertionOk := existingLiveDevice.(*panasonicNetCam.PanasonicNetworkCamera)
		if assertionOk {
			err = panNetCam.PatchAndReload(&panasonicNetCam.PanasonicNetworkCameraPatchParams{
				devicePE.Name, devicePE.State, devicePE.Brand,
				devicePE.Protocol, devicePE.Hostname, devicePE.Port, devicePE.Username, decryptedPassword,
				devicePE.FallbackRTSPPort,
			})
		} else {
			// warning ?
		}

	default:
		bgev.Messages.AddError(DVC_ERR_11299.NewMessage(devicePE.Code, devicePE.Brand))
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return
	}

	if err != nil {
		bgev.Messages.AddError(DVC_ERR_11202.NewMessage(devicePE.Code, err.Error()))
		localEv.LocalEventHub().PublishBgEvent(bgev)
		return
	}
}
