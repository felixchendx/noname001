package netcam

import (
	baseTyping "noname001/app/base/typing"

	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

// TODO: periodic cache invalidation / update
type t_cache struct {
	pDat  baseTyping.BaseDevicePersistenceData
	lDat  baseTyping.BaseDeviceLiveData
	opCap baseTyping.BaseDeviceOpCap
	hwDat baseTyping.BaseDeviceHardwareData

	rtspPort string

	// api call full caches
	deviceInfoWrapper *response.DeviceInfoWrapper
	deviceInfoAceI    apicall.APICallEventIntface

	streamInfos map[string]*liveBase.StreamInfo

	// temp: struct for errors ?
	tempErrDetails map[string]string
}

func (dev *PanasonicNetworkCamera) newCache() {
	dev.cache = &t_cache{
		streamInfos: make(map[string]*liveBase.StreamInfo),

		tempErrDetails: make(map[string]string),
	}

	dev.cache.pDat = baseTyping.BaseDevicePersistenceData{
		ID   : dev.id,
		Code : dev.code,
		Name : dev.name,
		State: dev.state,
		Brand: dev.brand,
	}
}
