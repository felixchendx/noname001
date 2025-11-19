package hikvision

import (
	appConstant "noname001/app/constant"

	"noname001/corebiz/integration/hikvision/httpapi"
	hikTyping "noname001/corebiz/integration/hikvision/httpapi/v1/typing"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

// func (dev *HikvisionDevice) GetDeviceInfo() (*liveBase.DeviceInfo, liveBase.DeviceResponseIntface) {
// 	var (
// 		deviceInfo *liveBase.DeviceInfo

// 		wrapped = dev.cache.deviceInfoWrapper
// 		aceI    = dev.cache.deviceInfoAceI
// 	)

// 	if wrapped == nil {
// 		wrapped, aceI = dev.api.FetchDeviceInfo(&httpapi.FetchDeviceInfoParams{})
// 	}

// 	if !aceI.IsConsideredError() {
// 		deviceInfo = wrapped.DeviceInfo
// 	}

// 	dev.cache.deviceInfoWrapper = wrapped
// 	dev.cache.deviceInfoAceI    = aceI

// 	return deviceInfo, aceI
// }

// TODO: cache
func (dev *HikvisionDevice) GetStreamInfo(channelID string, streamType appConstant.BrandStreamType) (*liveBase.StreamInfo, liveBase.DeviceResponseIntface) {
	var (
		streamInfo *liveBase.StreamInfo
	)

	rtspPortToUse := dev.cache.rtspPort
	if rtspPortToUse == "" { rtspPortToUse = dev.fallbackRTSPPort }

	wrapped, aceI := dev.api.FetchStreamInfo(&httpapi.FetchStreamInfoParams{
		ChannelID: channelID,
		StreamType: dev.mapStreamType(streamType),

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

func (dev *HikvisionDevice) mapStreamType(bst appConstant.BrandStreamType) (hikTyping.StreamType) {
	switch bst {
	case appConstant.HIKVISION__MAIN_STREAM : return hikTyping.MAIN_STREAM
	case appConstant.HIKVISION__SUB_STREAM  : return hikTyping.SUB_STREAM
	case appConstant.HIKVISION__THIRD_STREAM: return hikTyping.THIRD_STREAM
	}

	return hikTyping.MAIN_STREAM
}
