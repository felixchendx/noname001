package v1

type TXT_ResponseStatus struct {
	RequestURL string
	StatusCode int
	StatusMsg  string
}

type TXT_SoftwareVersion struct {
	Version     string
	ReleaseDate string
}

type TXT_SystemInfo struct {
	Processor    string
	SerialNumber string
	UpdateSerial string
}

type TXT_DeviceInfo struct {
	DeviceName      string
	DeviceType      string
	Manufacturer    string
	HardwareVersion string
	SoftwareVersion *TXT_SoftwareVersion
	SystemInfo      *TXT_SystemInfo
}

type TXT_VideoEncodeConfig struct {
	ChannelID           string
	ChannelRecordType   string // TODO: review this

	AudioSource         string
	AudioBitrate        int 	// in kbit/s
	AudioCompression    string
	AudioDepth          int
	AudioFrequency      int
	AudioMode           int
	AudioPack           string
	AudioPacketPeriod   int
	AudioEnable         bool

	VideoResolution     string
	VideoWidth          int
	VideoHeight         int
	VideoBitrate        int 	// in kbit/s
	VideoBitrateControl string
	VideoCompression    string
	VideoFps            int
	VideoGop            int
	VideoPack           string
	VideoProfile        string
	VideoQuality        int
	VideoQualityRange   int
	VideoEnable         bool
}
