package dahua

import (
	"time"
	
	baseTyping "noname001/app/base/typing"

	"noname001/corebiz/integration/dahua/httpapi"
)

// these are blocking calls

func (dev *DahuaDevice) capTest__fetchDeviceInfo() {
	wrapped, aceI := dev.api.FetchDeviceInfo(&httpapi.FetchDeviceInfoParams{})
	if aceI.IsConsideredError() {
		var errMsg string = ""

		switch {
		case aceI.IsGoError(): errMsg  = aceI.GoError().Error()
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

func (dev *DahuaDevice) capTest__fetchAnalogChannelsInfo() {
	wrapped, aceI := dev.api.FetchAnalogInputChannelList(&httpapi.FetchAnalogInputChannelListParams{})
	if aceI.IsConsideredError() {
		var errMsg string = ""

		switch {
		case aceI.IsGoError(): errMsg  = aceI.GoError().Error()
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

func (dev *DahuaDevice) capTest__fetchStreamInfo(channelID, streamType string) {
	wrapped, aceI := dev.api.FetchStreamInfo(&httpapi.FetchStreamInfoParams{
		ChannelID: channelID,
		// StreamType: streamType,
		RTSPPort: dev.fallbackRTSPPort,
	})
	if aceI.IsConsideredError() {
		var errMsg string = ""

		switch {
		case aceI.IsGoError(): errMsg  = aceI.GoError().Error()
		case aceI.IsAPIError(): errMsg = aceI.APIError().SimpleError()
		}

		dev.logger.Errorf("%s: ReadStreamInfo, %s", dev.logPrefix, errMsg)
		dev.cache.tempErrDetails["ReadStreamInfo"] = errMsg

	} else {
		dev.cache.opCap.CanReadStreamInfo = true

		dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_ALIVE
		dev.cache.lDat.LastSeen  = time.Now()
	}

	// streamID := fmt.Sprintf("%s%s", channelID, streamType)
	dev.cache.streamInfoWrappers[channelID] = wrapped
	dev.cache.streamInfoAceIs[channelID]    = aceI
}
