package coordinator

import (
	"fmt"
)

const (
	// [0] device_code
	// [1] channel_id
	// [2] stream_type
	mod_device__DEVICE_CHANNEL_PREVIEW string = "mod-device/device-channel-preview/%s/%s/%s"
)

func (coord *Coordinator) PublishDeviceChannelPreviewURL(deviceCode, channelID, streamType string) (string) {
	return coord.mediaServer.LocalRTSPPublisherBaseURL() +
		"/" + fmt.Sprintf(mod_device__DEVICE_CHANNEL_PREVIEW, deviceCode, channelID, streamType)
}

func (coord *Coordinator) DeviceChannelPreviewURL(requesterHostname, deviceCode, channelID, streamType, streamProtocol string) (string) {
	if (requesterHostname == "" || deviceCode == "" || channelID == "" || streamType == "") { return ""	}

	streamingPorts := coord.mediaServer.StreamingPorts()
	protocolPort, ok := streamingPorts[streamProtocol]
	if !ok { return "" }

	var streamURL string
	authnPair := coord.mediaServer.ViewerAuthnPair()

	switch streamProtocol {
	case "rtsp":
		streamURL = fmt.Sprintf("rtsp://%s@%s:%s", authnPair, requesterHostname, protocolPort) +
			"/" + fmt.Sprintf(mod_device__DEVICE_CHANNEL_PREVIEW, deviceCode, channelID, streamType)

	case "hls":
		streamURL = fmt.Sprintf("http://%s:%s", requesterHostname, protocolPort) +
			"/" + fmt.Sprintf(mod_device__DEVICE_CHANNEL_PREVIEW, deviceCode, channelID, streamType) +
			"/stream.m3u8"

	default:
		streamURL = ""
	}

	return streamURL
}
