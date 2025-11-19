package response

type AnalogInputChannelListWrapper struct {
	AnalogInputChannelList []*AnalogInputChannel

	OriginalData    map[string]any
}
type AnalogInputChannel struct {
	ID          string
	Enabled     bool
	Name        string
	VideoFormat string // PAL, NTSC
	PortType    string // SDI, OPT, VGA, HDMI, YPbPr
	ResDesc     string
}
// ^^^ tweak after reviewing api from other brands

type DigitalInputChannelListWrapper struct {
	DigitalInputChannelList []*DigitalInputChannel

	OriginalData    map[string]any
}
type DigitalInputChannel struct {
	ID                     string
	ProxyProtocol          string   // HIKVISION, SONY, ISAPI, ONVIF, ...
	Online                 bool
	StreamingChannelIDList []string // mostly []string{"X01", "X02"} where X = ChannelID
	ChannelDetectionResult string
}
// ^^^ tweak after reviewing api from other brands
