package live

import (
	baseTyping "noname001/app/base/typing"
)

// these stuffs are remains of old mechanism, revisit later

// =========================== VVV input def VVV ============================ //
type inputIntface interface {} 

type inputSourceModDevice struct {
	StreamURL string

	DeviceSnapshot   *baseTyping.BaseDeviceSnapshot
	DeviceStreamInfo *baseTyping.BaseDeviceStreamInfo
}

type inputSourceExternal struct {
	URL string
}

type inputSourceFile struct {
	Filepath string // abs filepath
}
// =========================== ^^^ input def ^^^ ============================ //

// =========================== VVV output def VVV =========================== //
type outputDestinationInternalMediaServer struct {
	PublishURL string

	TargetOutput *targetOutput
}

type targetOutput struct {
	StreamProfileCode string
	StreamProfileName string

	VideoFPS         float64
	VideoWidth       int
	VideoHeight      int
	VideoCodec       string
	VideoCompression int
	VideoBitrate     int

	AudioCodec       string
	AudioCompression int
	AudioBitrate     int

	ShowTimestamp    string
	ShowVideoInfo    string
	ShowAudioInfo    string
	ShowSiteInfo     string
}
// =========================== ^^^ output def ^^^ =========================== //
