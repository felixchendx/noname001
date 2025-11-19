package response

type StreamInfoWrapper struct {
	StreamURL    string
	StreamInfo   *StreamInfo

	OriginalData map[string]any
}

type StreamInfo struct {
	ChannelID   string
	ChannelName string
	Enabled bool
	// StreamType ? StreamID ? sub / main / 03

	VideoEnabled bool
	VideoCodecType string
	VideoResolutionWidth int
	VideoResolutionHeight int
	VideoFPS float32
	VideoBitrate int // in bit/s

	AudioEnabled bool
	AudioCodecType string // compression type ?
}
