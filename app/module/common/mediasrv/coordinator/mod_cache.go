package coordinator

import (
	"fmt"
)

const (
	mod_cache__SOURCE_STREAM_VIEW string = mod_stream__STREAM_VIEW

	// [0] node_id
	// [1] stream_item_code
	mod_cache__RELAYED_STREAM_VIEW string = "mod-cache/relayed-stream-view/%s/%s"
)

func (coord *Coordinator) SourceStreamURL(mediasrvAuthnPair, mediasrvIP, rtspPort, streamCode string) (string) {
	return fmt.Sprintf("rtsp://%s@%s:%s", mediasrvAuthnPair, mediasrvIP, rtspPort) +
		"/" + fmt.Sprintf(mod_cache__SOURCE_STREAM_VIEW, streamCode)
}

func (coord *Coordinator) RelayedStreamURI(nodeID, streamCode string) (string) {
	return fmt.Sprintf(mod_cache__RELAYED_STREAM_VIEW, nodeID, streamCode)
}

func (coord *Coordinator) RelayedStreamViewURL(requesterHostname, nodeID, streamCode, streamProtocol string) (string) {
	if (requesterHostname == "" || nodeID == "" || streamCode == "") { return "" }

	streamingPorts := coord.mediaServer.StreamingPorts()
	protocolPort, ok := streamingPorts[streamProtocol]
	if !ok { return "" }

	var streamURL string
	authnPair := coord.mediaServer.RelayAuthnPair()

	switch streamProtocol {
	case "rtsp":
		streamURL = fmt.Sprintf("rtsp://%s@%s:%s", authnPair, requesterHostname, protocolPort) +
			"/" + fmt.Sprintf(mod_cache__RELAYED_STREAM_VIEW, nodeID, streamCode)

	case "hls":
		// chrome based browser no longer support authn in URL (deprecated for 7years+)
		// and will strip authn in URL
		// so... authn is moved to request header
		streamURL = fmt.Sprintf("http://%s:%s", requesterHostname, protocolPort) +
			"/" + fmt.Sprintf(mod_cache__RELAYED_STREAM_VIEW, nodeID, streamCode) +
			"/stream.m3u8"

	default:
		streamURL = ""
	}

	return streamURL
}
