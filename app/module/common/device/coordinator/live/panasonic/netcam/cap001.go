package netcam

import (
	"fmt"
	"time"

	baseTyping  "noname001/app/base/typing"
	appConstant "noname001/app/constant"

	baseResponse "noname001/corebiz/integration/base/response"

	camInterface "noname001/corebiz/integration/panasonic/netcam"
	camV1Typing  "noname001/corebiz/integration/panasonic/netcam/v1/typing"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

func (dev *PanasonicNetworkCamera) capTest__fetchDeviceInfo() {
	wrapped, aceI := dev.api.FetchDeviceInfo(&camInterface.FetchDeviceInfoParams{})
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

func (dev *PanasonicNetworkCamera) capTest__fetchDigitalChannelsInfo() {
	// all other informations also available on cached deviceInfo
	wrapped := dev.cache.deviceInfoWrapper

	if wrapped != nil && wrapped.OriginalData != nil {
		untypedProductInfo, hasOriProductInfo := wrapped.OriginalData["productInfo"]
		if hasOriProductInfo {
			typedProductInfo, assertionOK := untypedProductInfo.(*camV1Typing.ProductInformation)
			if assertionOK {
				dev.cache.opCap.CanReadDigitalInputChannels = true

				_digitalChannels := make([]baseTyping.BaseDeviceDigitalChannel, 0)

				_diChan1 := baseTyping.BaseDeviceDigitalChannel{
					ChannelID  : "1",
					ChannelName: "Ch 1",
					Enabled    : false,
				}
				if typedProductInfo.ITransmit_h264 == "1" { _diChan1.Enabled = true }
				_digitalChannels = append(_digitalChannels, _diChan1)

				dev.cache.hwDat.LastUpdated     = time.Now()
				dev.cache.hwDat.DigitalChannels = _digitalChannels
			}
		}
	}
}

func (dev *PanasonicNetworkCamera) capTest__fetchStreamInfo() {
	// all other informations also available on cached deviceInfo
	wrapped := dev.cache.deviceInfoWrapper

	if wrapped != nil && wrapped.OriginalData != nil {
		untypedProductInfo, hasOriProductInfo := wrapped.OriginalData["productInfo"]
		if hasOriProductInfo {
			typedProductInfo, assertionOK := untypedProductInfo.(*camV1Typing.ProductInformation)
			if assertionOK {
				dev.cache.opCap.CanReadStreamInfo = true

				{
					streamType := appConstant.PANASONIC_NETCAM__STREAM_1
					streamKey  := fmt.Sprintf("1-%s", streamType)

					_baseStreamInfo := &baseResponse.StreamInfo{
						ChannelID  : streamKey,
						ChannelName: "Ch 1 - Stream(1)",
						Enabled    : false,

						VideoEnabled         : false,
						VideoCodecType       : "",
						VideoResolutionWidth : 0,
						VideoResolutionHeight: 0,
						VideoFPS             : 0,
						VideoBitrate         : 0, // in bit/s

						AudioEnabled  : false,
						AudioCodecType: "",
					}
					if typedProductInfo.ITransmit_h264 == "1" {
						_baseStreamInfo.Enabled = true
						_baseStreamInfo.VideoEnabled = true
					}
					_baseStreamInfo.VideoCodecType = camV1Typing.TranslateStreamEncode(typedProductInfo.StreamEncode)
					_baseStreamInfo.VideoResolutionWidth, _baseStreamInfo.VideoResolutionHeight = camV1Typing.TranslateResolution(
						typedProductInfo.Ratio, typedProductInfo.IResolution_h264,
					)
					_baseStreamInfo.VideoBitrate = camV1Typing.TranslateBitrate(typedProductInfo.IBitrate_h264)

					dev.cache.streamInfos[streamKey] = &liveBase.StreamInfo{
						_baseStreamInfo,
						dev._determineStreamURL(streamType),
					}
				}
				{
					streamType := appConstant.PANASONIC_NETCAM__STREAM_2
					streamKey  := fmt.Sprintf("1-%s", streamType)

					_baseStreamInfo := &baseResponse.StreamInfo{
						ChannelID  : streamKey,
						ChannelName: "Ch 1 - Stream(2)",
						Enabled    : false,

						VideoEnabled         : false,
						VideoCodecType       : "",
						VideoResolutionWidth : 0,
						VideoResolutionHeight: 0,
						VideoFPS             : 0,
						VideoBitrate         : 0, // in bit/s

						AudioEnabled  : false,
						AudioCodecType: "",
					}
					if typedProductInfo.ITransmit_h264_2 == "1" {
						_baseStreamInfo.Enabled = true
						_baseStreamInfo.VideoEnabled = true
					}
					_baseStreamInfo.VideoCodecType = camV1Typing.TranslateStreamEncode(typedProductInfo.StreamEncode_2)
					_baseStreamInfo.VideoResolutionWidth, _baseStreamInfo.VideoResolutionHeight = camV1Typing.TranslateResolution(
						typedProductInfo.Ratio, typedProductInfo.IResolution_h264_2,
					)
					_baseStreamInfo.VideoBitrate = camV1Typing.TranslateBitrate(typedProductInfo.IBitrate_h264_2)

					dev.cache.streamInfos[streamKey] = &liveBase.StreamInfo{
						_baseStreamInfo,
						dev._determineStreamURL(streamType),
					}
				}
				{
					streamType := appConstant.PANASONIC_NETCAM__STREAM_3
					streamKey  := fmt.Sprintf("1-%s", streamType)

					_baseStreamInfo := &baseResponse.StreamInfo{
						ChannelID  : streamKey,
						ChannelName: "Ch 1 - Stream(3)",
						Enabled    : false,

						VideoEnabled         : false,
						VideoCodecType       : "",
						VideoResolutionWidth : 0,
						VideoResolutionHeight: 0,
						VideoFPS             : 0,
						VideoBitrate         : 0, // in bit/s

						AudioEnabled  : false,
						AudioCodecType: "",
					}
					if typedProductInfo.ITransmit_h264_3 == "1" {
						_baseStreamInfo.Enabled = true
						_baseStreamInfo.VideoEnabled = true
					}
					_baseStreamInfo.VideoCodecType = camV1Typing.TranslateStreamEncode(typedProductInfo.StreamEncode_3)
					_baseStreamInfo.VideoResolutionWidth, _baseStreamInfo.VideoResolutionHeight = camV1Typing.TranslateResolution(
						typedProductInfo.Ratio, typedProductInfo.IResolution_h264_3,
					)
					_baseStreamInfo.VideoBitrate = camV1Typing.TranslateBitrate(typedProductInfo.IBitrate_h264_3)

					dev.cache.streamInfos[streamKey] = &liveBase.StreamInfo{
						_baseStreamInfo,
						dev._determineStreamURL(streamType),
					}
				}
				{
					streamType := appConstant.PANASONIC_NETCAM__STREAM_4
					streamKey  := fmt.Sprintf("1-%s", streamType)

					_baseStreamInfo := &baseResponse.StreamInfo{
						ChannelID  : streamKey,
						ChannelName: "Ch 1 - Stream(4)",
						Enabled    : false,

						VideoEnabled         : false,
						VideoCodecType       : "",
						VideoResolutionWidth : 0,
						VideoResolutionHeight: 0,
						VideoFPS             : 0,
						VideoBitrate         : 0, // in bit/s

						AudioEnabled  : false,
						AudioCodecType: "",
					}
					if typedProductInfo.ITransmit_h264_4 == "1" {
						_baseStreamInfo.Enabled = true
						_baseStreamInfo.VideoEnabled = true
					}
					_baseStreamInfo.VideoCodecType = camV1Typing.TranslateStreamEncode(typedProductInfo.StreamEncode_4)
					_baseStreamInfo.VideoResolutionWidth, _baseStreamInfo.VideoResolutionHeight = camV1Typing.TranslateResolution(
						typedProductInfo.Ratio, typedProductInfo.IResolution_h264_4,
					)
					_baseStreamInfo.VideoBitrate = camV1Typing.TranslateBitrate(typedProductInfo.IBitrate_h264_4)

					dev.cache.streamInfos[streamKey] = &liveBase.StreamInfo{
						_baseStreamInfo,
						dev._determineStreamURL(streamType),
					}
				}

				dev.cache.hwDat.LastUpdated = time.Now()
			}
		}
	}
}

// temp, move to somewhere within package noname001/corebiz/integration/panasonic/netcam
func (dev *PanasonicNetworkCamera) _determineStreamURL(streamType appConstant.BrandStreamType) (string) {
	rtspPortToUse := dev.cache.rtspPort
	if rtspPortToUse == "" { rtspPortToUse = dev.fallbackRTSPPort }

	streamURI := ""
	streamURL := ""

	switch streamType {
	case appConstant.PANASONIC_NETCAM__STREAM_1: streamURI = "stream_1"
	case appConstant.PANASONIC_NETCAM__STREAM_2: streamURI = "stream_2"
	case appConstant.PANASONIC_NETCAM__STREAM_3: streamURI = "stream_3"
	case appConstant.PANASONIC_NETCAM__STREAM_4: streamURI = "stream_4"

	default:
		streamURI = "stream_1"
	}

	if streamURI == "" {
		streamURL = ""

	} else {
		// 2.3.3. RTSP URL
		streamURL = fmt.Sprintf(
			"rtsp://%s:%s@%s:%s/Src/MediaInput/%s",
			dev.username, dev.password,
			dev.hostname, rtspPortToUse,
			streamURI,
		)
	}

	return streamURL
}
