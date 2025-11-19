package dahua

import (
	baseTyping "noname001/app/base/typing"

	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"
)

type t_cache struct {
	pDat  baseTyping.BaseDevicePersistenceData
	lDat  baseTyping.BaseDeviceLiveData
	opCap baseTyping.BaseDeviceOpCap
	hwDat baseTyping.BaseDeviceHardwareData

	rtspPort string

	// api call full caches
	deviceInfoWrapper *response.DeviceInfoWrapper
	deviceInfoAceI    apicall.APICallEventIntface

	analogInputChannelListWrapper *response.AnalogInputChannelListWrapper
	analogInputChannelListAceI    apicall.APICallEventIntface

	digitalInputChannelListWrapper *response.DigitalInputChannelListWrapper
	digitalInputChannelListAceI    apicall.APICallEventIntface

	streamInfoWrappers map[string]*response.StreamInfoWrapper
	streamInfoAceIs    map[string]apicall.APICallEventIntface

	// temp: struct for errors ?
	tempErrDetails map[string]string
}

func (dev *DahuaDevice) newCache() {
	dev.cache = &t_cache{
		streamInfoWrappers: make(map[string]*response.StreamInfoWrapper),
		streamInfoAceIs   : make(map[string]apicall.APICallEventIntface),

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
