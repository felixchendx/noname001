package v1

import (
	"encoding/xml"
)

// last tested: 2024 04 05
// these are tested against development setup
// which returns Version 1.0

// TODO(FUTURE): add more firmware info / sdk info

// all the "xx.x.xxx XML_xxx" comment above type
// corresponds to file ISAPI__Version_2_0__Sept_2019.pdf unless stated otherwise
// reference with precaution (v1.0 vs v2.0)

// 16.2.253 XML_ResponseStatus
type XML_ResponseStatus struct {
	XMLName       xml.Name `xml:"ResponseStatus"`
	Version       string   `xml:"version,attr"`
	Xmlns         string   `xml:"xmlns,attr"`

	RequestURL    string            `xml:"requestURL"`
	StatusCode    int               `xml:"statusCode"`
	StatusString  string            `xml:"statusString"`
	ID            int               `xml:"id"`
	SubStatusCode string            `xml:"subStatusCode"`
	ErrorCode     int               `xml:"errorCode"`
	ErrorMsg      string            `xml:"errorMsg"`
	AdditionalErr XML_AdditionalErr `xml:"AdditionalErr"`
}
type XML_AdditionalErr struct {
	StatusList []XML_Status `xml:"StatusList>Status"`
}
type XML_Status struct {
	ID            string `xml:"id"`
	StatusCode    int    `xml:"statusCode"`
	StatusString  string `xml:"statusString"`
	SubStatusCode string `xml:"subStatusCode"`
}

// 16.2.115 XML_DeviceInfo
type XML_DeviceInfo struct {
	XMLName xml.Name `xml:"DeviceInfo"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	DeviceName           string `xml:"deviceName"`
	DeviceID             string `xml:"deviceID"`
	Model                string `xml:"model"`
	SerialNumber         string `xml:"serialNumber"`
	MACAddress           string `xml:"macAddress"`
	FirmwareVersion      string `xml:"firmwareVersion"`
	FirmwareReleasedDate string `xml:"firmwareReleasedDate"`
	EncoderVersion       string `xml:"encoderVersion"`
	EncoderReleasedDate  string `xml:"encoderReleasedDate"`
	DeviceType           string `xml:"deviceType"`
	TelecontrolID        int    `xml:"telecontrolID"`
	HardwareVersion      string `xml:"hardwareVersion"`
	Manufacturer         string `xml:"manufacturer"`
	CustomizedInfo       string `xml:"customizedInfo"`
}

// 16.2.280 XML_StreamingChannelList
type XML_StreamingChannelList struct {
	XMLName xml.Name `xml:"StreamingChannelList"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	StreamingChannels []XML_StreamingChannel `xml:"StreamingChannel"`
}
// 16.2.279 XML_StreamingChannel
type XML_StreamingChannel struct {
	XMLName xml.Name `xml:"StreamingChannel"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	ID          string `xml:"id"`
	ChannelName string `xml:"channelName"`
	Enabled     bool   `xml:"enabled"`

	Transport XML_SC_Transport `xml:"Transport"`
	Video     XML_SC_Video     `xml:"Video"`
	Audio     XML_SC_Audio     `xml:"Audio"`

	EnableCABAC        bool `xml:"enableCABAC"`
	SubStreamRecStatus bool `xml:"subStreamRecStatus"`
	CustomStreamEnable bool `xml:"customStreamEnable"`
}
// NOTE: there are many more fields not included in this struct, add as necessary
type XML_SC_Transport struct {
	ControlProtocolList []XML_SC_ControlProtocol `xml:"ControlProtocolList"`
}
type XML_SC_ControlProtocol struct {
	StreamingTransport string `xml:"streamingTransport"`
}
type XML_SC_Video struct {
	Enabled             bool   `xml:"enabled"`
	VideoInputChannelID string `xml:"videoInputChannelID"`

	VideoCodecType        string `xml:"videoCodecType"` // MPEG4,MJPEG,3GP,H.264,HK.264,MPNG,SVAC,H.265
	// VideoScanType ???? not in doc
	VideoResolutionWidth  int    `xml:"videoResolutionWidth"`
	VideoResolutionHeight int    `xml:"videoResolutionHeight"`
	
	VideoQualityControlType string `xml:"videoQualityControlType"` // CBR,VBR
	ConstantBitRate int `xml:"constantBitRate"` // in kbps
	// FixedQuality ???? not in doc
	VBRUpperCap     int `xml:"vbrUpperCap"`     // in kbps
	VBRLowerCap     int `xml:"vbrLowerCap"`     // in kbps
	MaxFrameRate    int `xml:"maxFrameRate"`    // value is multiplied by 100

	KeyFrameInterval int  `xml:"keyFrameInterval"` // in milliseconds
	RotationDegree   int  `xml:"RotationDegree"`   // degress 0..360
	MirrorEnabled    bool `xml:"mirrorEnabled"`

	SnapShotImageType string `xml:"snapShotImageType"` // JPEG,GIF,PNG

	MPEG4Profile string `xml:"Mpeg4Profile"` // SP, ASP
	H264Profile  string `xml:"H264Profile"`  // Baseline,Main,High, Extended
	SVACProfile  string `xml:"SVACProfile"`  // Baseline,Main,High,Extended
	
	GovLength int `xml:"GovLength"`
	// SVC 
	Smoothing int `xml:"smoothing"`
	// SmartCodec

	VBRAverageCap int `xml:"vbrAverageCap"` // in kbps
	// IntelligentInfoDisplayMethod
}
type XML_SC_Audio struct {
	Enabled             bool   `xml:"enabled"`
	AudioInputChannelID string `xml:"audioInputChannelID"`
	
	AudioCompressionType        string `xml:"audioCompressionType"`        // G.711alaw,G.711ulaw,G.726,G.729,G.729a,G.729b,PCM,MP3,AC3,AAC,ADPCM,MP2L2
	AudioInboundCompressionType string `xml:"audioInboundCompressionType"` // G.711alaw,G.711ulaw,G.726,G.729,G.729a,G.729b,PCM,MP3,AC3,AAC,ADPCM,MP2L2

	AudioBitRate      int     `xml:"audioBitRate"`      // in kbps
	AudioSamplingRate float64 `xml:"audioSamplingRate"` // in kHz
	AudioResolution   int     `xml:"audioResolution"`   // in bits

	// VoiceChanger
}

// 16.2.321 XML_VideoInput
type XML_VideoInput struct {
	XMLName xml.Name `xml:"VideoInput"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	VideoInputChannelList []XML_VideoInputChannel `xml:"VideoInputChannelList>VideoInputChannel"`
}
// 16.2.322 XML_VideoInputChannel
type XML_VideoInputChannel struct {
	XMLName xml.Name `xml:"VideoInputChannel"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	ID                string `xml:"id"`
	InputPort         string `xml:"inputPort"`
	VideoInputEnabled bool   `xml:"videoInputEnabled"`
	Name              string `xml:"name"`
	VideoFormat       string `xml:"videoFormat"` // PAL, NTSC
	PortType          string `xml:"portType"`    // SDI, OPT, VGA, HDMI, YPbPr
	ResDesc           string `xml:"resDesc"`
}

// 16.2.173 XML_InputProxyChannelStatusList
type XML_InputProxyChannelStatusList struct {
	XMLName xml.Name `xml:"InputProxyChannelStatusList"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	InputProxyChannelStatusList []XML_InputProxyChannelStatus `xml:"InputProxyChannelStatus"`
}
// 16.2.172 XML_InputProxyChannelStatus
type XML_InputProxyChannelStatus struct {
	XMLName xml.Name `xml:"InputProxyChannelStatus"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	ID string `xml:"id"`

	// WARNING! this struct's field differs from testing hardware vs doc
	// also does not has specific structure in doc ? 
	SourceInputPortDescriptor XML_IPCS_sourceInputPortDescriptor `xml:"sourceInputPortDescriptor"`

	Online bool `xml:"online"`

	StreamingProxyChannelIDList []XML_IPCS_streamingProxyChannelId `xml:"streamingProxyChannelIdList"`
	// RelatedIOProxy ???

	ChanDetectResult string `xml:"chanDetectResult"`
	// possible values from doc...
	// network camera status: "connect"-connected, "overSysBandwidth"-insufficient bandwidth,
	// "domainError"-incorrect domain name, "ipcStreamFail"-getting stream failed, "connecting", "chacnNoError"-incorrect
	// channel No., "cipAddrConflictWithDev": IP address is conflicted with device address, "ipAddrConflicWithIpc"-IP
	// address conflicted, "errorUserNameOrPasswd"-incorrect user name or password, "netUnreachable"-invalid network
	// address, "unknownError"-unknown error, "notExist"-does not exist, "ipcStreamTypeNotSupport"-the stream
	// transmission mode is not supported, "ipcResolutionNotSupport"-the resolution of network camera is not supported

	// SecurityStatus ???
}
// 16.2.276 XML_sourceDescriptor
// ^^^ the only descriptor with specific number and doc
// VVV differs from the descriptor above, also has no specific number...
type XML_IPCS_sourceInputPortDescriptor struct {
	XMLName xml.Name `xml:"sourceInputPortDescriptor"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`

	// === VVV FIELD FROM DOCS VVV ===
	AdminProtocol string `xml:"adminProtocol"` // HIKVISION, SONY, ISAPI, ONVIF, ...
	// === ^^^ FIELD FROM DOCS ^^^ ===

	// === VVV FIELD FROM DEV HW VVV ===
	ProxyProtocol string `xml:"proxyProtocol"` // HIKVISION, ONVIF, ...
	// === ^^^ FIELD FROM DEV HW ^^^ ===

	AddressingFormatType string `xml:"addressingFormatType"` // ipaddress, hostname
	Hostname             string `xml:"hostName"`
	IPAddress            string `xml:"ipAddress"`
	IPv6Address          string `xml:"ipv6Address"`
	ManagePortNo         int    `xml:"managePortNo"`
	SrcInputPort         string `xml:"srcInputPort"`
	Username             string `xml:userName`
	Password             string `xml:password`
	StreamType           string `xml:"streamType"` // auto, tcp, udp
	DeviceID             string `xml:"deviceID"`
	DeviceTypeName       string `xml:"deviceTypeName"`
	SerialNumber         string `xml:"serialNumber"`
	FirmwareVersion      string `xml:"firmwareVersion"`
	FirmwareCode         string `xml:"firmwareCode"`
}
type XML_IPCS_streamingProxyChannelId struct {
	ID string `xml:"streamingProxyChannelId"`
}