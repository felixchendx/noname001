package util

import (
	"fmt"
)

func GenerateStreamURL(hostname, port, user, pass, streamID string, streamType string) (string) {
	streamURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%s/cam/realmonitor?channel=%s&subtype=%s",
		user, pass,
		hostname, port,
		streamID, 
		streamType,
	)
	return streamURL
}
