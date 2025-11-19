package stream

type StreamSnapshot struct {
	ID   string
	Code string

	Persistence StreamPersistenceData
	Live        StreamLiveData

	// === source type specific data ===
	// ModDevStreamInfo *ModDeviceStreamInfo
}

type StreamPersistenceData struct {
	// === item ===
	State            string
	SourceType       string // mod_device / external / file
	DeviceCode       string
	DeviceChannelID  string
	DeviceStreamType string
	ExternalURL      string
	Filepath         string
	EmbeddedFilepath string

	// === group - temp ===
	GroupID string
	// GroupCode string
	GroupState string

	// === profile ===
	ProfileID   string
	ProfileCode string
	ProfileName string

	VideoEnabled           bool
	TargetVideoCodec       string
	TargetVideoCompression int
	TargetVideoBitrate     int

	AudioEnabled           bool
	TargetAudioCodec       string
	TargetAudioCompression int
	TargetAudioBitrate     int
}

type StreamLiveData struct {
	State         LiveState
	FailState     LiveFailState
	StreamerState LiveStreamerState

	URI string
	URL string

	// ffprobe / standard AV informations
	// InputProbes  FFProbeData
	// OutputProbes FFProbeData

	// temp
	EstimatedOutputVideoBitrate int
}

type FFProbeData struct {
	Streams []string
	Format  string
}


type ModDeviceStreamInfo struct {
	
	
	// DeviceID              string
	// DeviceName            string
	// DeviceType            string

	// ChannelID             string
	// ChannelName           string

	// VideoEnabled          bool
	// VideoResolutionWidth  int
	// VideoResolutionHeight int
	// VideoCodec            string
	// VideoFPS              float32
	// VideoBitrate          int

	// AudioEnabled          bool
	// AudioCodec            string
	// AudioBitrate          int
}
