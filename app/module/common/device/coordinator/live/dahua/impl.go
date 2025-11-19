package dahua

import (
	appConstant "noname001/app/constant"

	"noname001/corebiz/integration/dahua/httpapi"
	dahuaTyping "noname001/corebiz/integration/dahua/httpapi/v1/typing"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

// func (dev *DahuaDevice) GetDeviceInfo() (*liveBase.DeviceInfo, liveBase.DeviceResponseIntface) {
// 	var (
// 		deviceInfo *liveBase.DeviceInfo
// 	)

// 	wrapped, aceI := dev.api.FetchDeviceInfo(&httpapi.FetchDeviceInfoParams{})
// 	if !aceI.IsConsideredError() {
// 		deviceInfo = wrapped.DeviceInfo
// 	}

// 	return deviceInfo, aceI
// }

func (dev *DahuaDevice) GetStreamInfo(channelID string, streamType appConstant.BrandStreamType) (*liveBase.StreamInfo, liveBase.DeviceResponseIntface) {
	var (
		streamInfo *liveBase.StreamInfo
	)

	rtspPortToUse := dev.cache.rtspPort
	if rtspPortToUse == "" { rtspPortToUse = dev.fallbackRTSPPort }

	wrapped, aceI := dev.api.FetchStreamInfo(&httpapi.FetchStreamInfoParams{
		ChannelID: channelID,
		// StreamType: dev.mapStreamType(streamType),

		RTSPPort: rtspPortToUse,
	})
	if !aceI.IsConsideredError() {
		streamInfo = &liveBase.StreamInfo{
			wrapped.StreamInfo,
			wrapped.StreamURL,
		}
	}

	return streamInfo, aceI
}

func (dev *DahuaDevice) mapStreamType(bst appConstant.BrandStreamType) (dahuaTyping.StreamType) {
	switch bst {
	case appConstant.DAHUA__MAIN_STREAM   : return dahuaTyping.MAIN_STREAM
	case appConstant.DAHUA__EXTRA_STREAM_1: return dahuaTyping.EXTRA_STREAM_1
	case appConstant.DAHUA__EXTRA_STREAM_2: return dahuaTyping.EXTRA_STREAM_2
	case appConstant.DAHUA__EXTRA_STREAM_3: return dahuaTyping.EXTRA_STREAM_3
	}

	return dahuaTyping.MAIN_STREAM
}
