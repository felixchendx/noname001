package runner

type RawRunnerConfig struct {
	CfgDef               string               `yaml:"cfgdef"`
	Execute              bool                 `yaml:"execute"`

	SegmentStreamProfile SegmentStreamProfile `yaml:"stream_profile"`
	SegmentStreamGroup   SegmentStreamGroup   `yaml:"stream_group"`
}

type SegmentStreamProfile struct {
	Mode        string          `yaml:"mode"`

	Profiles    []StreamProfile `yaml:"profiles"`
}
type StreamProfile struct {
	Code                   string `yaml:"code"`
	Name                   string `yaml:"name"`
	State                  string `yaml:"state"`
	Note                   string `yaml:"note"`

	TargetVideoCodec       string `yaml:"target_video_codec"`
	TargetVideoCompression int    `yaml:"target_video_compression"`
	TargetVideoBitrate     int    `yaml:"target_video_bitrate"`

	TargetAudioCodec       string `yaml:"target_audio_codec"`
	TargetAudioCompression int    `yaml:"target_audio_compression"`
	TargetAudioBitrate     int    `yaml:"target_audio_bitrate"`
}

type SegmentStreamGroup struct {
	Mode   string        `yaml:"mode"`

	Groups []StreamGroup `yaml:"groups"`
}
type StreamGroup struct {
	Code              string       `yaml:"code"`
	Name              string       `yaml:"name"`
	State             string       `yaml:"state"`
	Note              string       `yaml:"note"`

	StreamProfileCode string       `yaml:"stream_profile_code"`

	Items             []StreamItem `yaml:"items"`
}
type StreamItem  struct {
	Code              string `yaml:"code"`
	Name              string `yaml:"name"`
	State             string `yaml:"state"`
	Note              string `yaml:"note"`

	SourceType       string `yaml:"source_type"`
	DeviceCode       string `yaml:"device_code"`
	DeviceChannelID  string `yaml:"device_channel_id"`
	ExternalURL      string `yaml:"external_url"`
	Filepath         string `yaml:"filepath"`
	EmbeddedFilepath string `yaml:"embedded_filepath"`
}
