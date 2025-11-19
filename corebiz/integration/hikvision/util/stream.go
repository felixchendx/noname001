package util

import (
	"fmt"
)

func GenerateStreamURL(hostname, port, user, pass string, channelID string, streamType string) (string) {
	streamURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%s/ISAPI/Streaming/channels/%s%s",
		user, pass,
		hostname, port,
		channelID, streamType,
	)
	return streamURL
}
