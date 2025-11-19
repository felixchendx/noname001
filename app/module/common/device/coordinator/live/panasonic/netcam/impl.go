package netcam

import (
	"fmt"

	appConstant "noname001/app/constant"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

func (dev *PanasonicNetworkCamera) GetStreamInfo(channelID string, streamType appConstant.BrandStreamType) (*liveBase.StreamInfo, liveBase.DeviceResponseIntface) {
	var (
		streamInfo *liveBase.StreamInfo
	)

	// TODO
	// temp, is called from live_stream even when device is not ready
	if dev.cache.deviceInfoAceI == nil {
		dev.capTest__fetchDeviceInfo()
	}

	if string(streamType) == "" {
		streamType = appConstant.PANASONIC_NETCAM__STREAM_1
	}

	streamKey := fmt.Sprintf("%s-%s", channelID, streamType)

	switch streamType {
	case appConstant.PANASONIC_NETCAM__STREAM_1: streamInfo, _ = dev.cache.streamInfos[streamKey]
	case appConstant.PANASONIC_NETCAM__STREAM_2: streamInfo, _ = dev.cache.streamInfos[streamKey]
	case appConstant.PANASONIC_NETCAM__STREAM_3: streamInfo, _ = dev.cache.streamInfos[streamKey]
	case appConstant.PANASONIC_NETCAM__STREAM_4: streamInfo, _ = dev.cache.streamInfos[streamKey]
	}

	return streamInfo, dev.cache.deviceInfoAceI
}
