package hikvision

import (
	"fmt"
	"time"

	baseTyping  "noname001/app/base/typing"
	appConstant "noname001/app/constant"

	"noname001/corebiz/integration/hikvision/httpapi"
)

// TODO: state guard for active -> inactive, when stuffs are still runnin

// these are blocking calls

func (dev *HikvisionDevice) capTest__fetchDeviceInfo() {
	wrapped, aceI := dev.api.FetchDeviceInfo(&httpapi.FetchDeviceInfoParams{})
	if aceI.IsConsideredError() {
		var errMsg string = ""

		switch {
		case aceI.IsGoError() : errMsg = aceI.GoError().Error()
		case aceI.IsAPIError(): errMsg = aceI.APIError().SimpleError()
		}

		dev.logger.Errorf("%s: ReadDeviceInfo, %s", dev.logPrefix, errMsg)
		dev.cache.tempErrDetails["ReadDeviceInfo"] = errMsg
		dev.cache.lDat.ConnStateMessage = errMsg

	} else {
		dev.cache.opCap.CanReadDeviceInfo = true

		dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_ALIVE
		dev.cache.lDat.LastSeen  = time.Now()

		dev.cache.hwDat.LastUpdated = time.Now()
		dev.cache.hwDat.DeviceID    = wrapped.DeviceInfo.DeviceID
		dev.cache.hwDat.DeviceName  = wrapped.DeviceInfo.DeviceName
		dev.cache.hwDat.Model       = wrapped.DeviceInfo.Model
		dev.cache.hwDat.DeviceType  = wrapped.DeviceInfo.DeviceType
	}

	dev.cache.deviceInfoWrapper = wrapped
	dev.cache.deviceInfoAceI    = aceI
}

func (dev *HikvisionDevice) capTest__fetchAnalogChannelsInfo() {
	wrapped, aceI := dev.api.FetchAnalogInputChannelList(&httpapi.FetchAnalogInputChannelListParams{})
	if aceI.IsConsideredError() {
		var errMsg string = ""

		switch {
		case aceI.IsGoError() : errMsg = aceI.GoError().Error()
		case aceI.IsAPIError(): errMsg = aceI.APIError().SimpleError()
		}

		dev.logger.Errorf("%s: ReadAnalogInputChannels, %s", dev.logPrefix, errMsg)
		dev.cache.tempErrDetails["ReadAnalogInputChannels"] = errMsg

	} else {
		dev.cache.opCap.CanReadAnalogInputChannels = true

		dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_ALIVE
		dev.cache.lDat.LastSeen  = time.Now()

		_analogChannels := make([]baseTyping.BaseDeviceAnalogChannel, 0)
		for _, _anChan := range wrapped.AnalogInputChannelList {
			_analogChannels = append(_analogChannels, baseTyping.BaseDeviceAnalogChannel{
				ChannelID  : _anChan.ID,
				ChannelName: _anChan.Name,
				Enabled    : _anChan.Enabled,
			})
		}
		dev.cache.hwDat.LastUpdated    = time.Now()
		dev.cache.hwDat.AnalogChannels = _analogChannels
	}

	dev.cache.analogInputChannelListWrapper = wrapped
	dev.cache.analogInputChannelListAceI    = aceI
}

func (dev *HikvisionDevice) capTest__fetchDigitalChannelsInfo() {
	wrapped, aceI := dev.api.FetchDigitalInputChannelList(&httpapi.FetchDigitalInputChannelListParams{})
	if aceI.IsConsideredError() {
		var errMsg string = ""

		switch {
		case aceI.IsGoError() : errMsg = aceI.GoError().Error()
		case aceI.IsAPIError(): errMsg = aceI.APIError().SimpleError()
		}

		dev.logger.Errorf("%s: ReadDigitalInputChannels, %s", dev.logPrefix, errMsg)
		dev.cache.tempErrDetails["ReadDigitalInputChannels"] = errMsg

	} else {
		dev.cache.opCap.CanReadDigitalInputChannels = true

		dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_ALIVE
		dev.cache.lDat.LastSeen  = time.Now()

		_digitalChannels := make([]baseTyping.BaseDeviceDigitalChannel, 0)
		for _, _diChan := range wrapped.DigitalInputChannelList {
			_digitalChannels = append(_digitalChannels, baseTyping.BaseDeviceDigitalChannel{
				ChannelID  : _diChan.ID,
				ChannelName: "-",
				Enabled    : _diChan.Online,
			})
		}
		dev.cache.hwDat.LastUpdated     = time.Now()
		dev.cache.hwDat.DigitalChannels = _digitalChannels
	}

	dev.cache.digitalInputChannelListWrapper = wrapped
	dev.cache.digitalInputChannelListAceI    = aceI
}

func (dev *HikvisionDevice) capTest__fetchStreamInfo(channelID string, streamType appConstant.BrandStreamType) {
	wrapped, aceI := dev.api.FetchStreamInfo(&httpapi.FetchStreamInfoParams{
		ChannelID: channelID,
		StreamType: dev.mapStreamType(streamType),
		RTSPPort: dev.fallbackRTSPPort,
	})
	if aceI.IsConsideredError() {
		var errMsg string = ""

		switch {
		case aceI.IsGoError() : errMsg = aceI.GoError().Error()
		case aceI.IsAPIError(): errMsg = aceI.APIError().SimpleError()
		}

		dev.logger.Errorf("%s: ReadStreamInfo, %s", dev.logPrefix, errMsg)
		dev.cache.tempErrDetails["ReadStreamInfo"] = errMsg

	} else {
		dev.cache.opCap.CanReadStreamInfo = true

		dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_ALIVE
		dev.cache.lDat.LastSeen  = time.Now()
	}

	streamKey := fmt.Sprintf("%s-%s", channelID, streamType)
	dev.cache.streamInfoWrappers[streamKey] = wrapped
	dev.cache.streamInfoAceIs[streamKey]    = aceI
}
