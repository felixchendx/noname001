package coordinator

import (
	"fmt"
)

const (
	// [0] stream_item_code
	mod_stream__STREAM_VIEW string = "mod-stream/stream-view/%s"
)


func (coord *Coordinator) PublishStreamViewURL(streamItemCode string) (string) {
	return coord.mediaServer.LocalRTSPPublisherBaseURL() +
		"/" + fmt.Sprintf(mod_stream__STREAM_VIEW, streamItemCode)
}

func (coord *Coordinator) StreamViewURL(requesterHostname, streamItemCode, streamProtocol string) (string) {
	if (requesterHostname == "" || streamItemCode == "") { return "" }

	streamingPorts := coord.mediaServer.StreamingPorts()
	protocolPort, ok := streamingPorts[streamProtocol]
	if !ok { return "" }

	var streamURL string
	authnPair := coord.mediaServer.ViewerAuthnPair()

	switch streamProtocol {
	case "rtsp":
		streamURL = fmt.Sprintf("rtsp://%s@%s:%s", authnPair, requesterHostname, protocolPort) +
			"/" + fmt.Sprintf(mod_stream__STREAM_VIEW, streamItemCode)

	case "hls":
		streamURL = fmt.Sprintf("http://%s:%s", requesterHostname, protocolPort) +
			"/" + fmt.Sprintf(mod_stream__STREAM_VIEW, streamItemCode) +
			"/stream.m3u8"

	default:
		streamURL = ""
	}

	return streamURL
}
