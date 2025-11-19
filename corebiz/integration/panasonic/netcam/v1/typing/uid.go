 package typing

// 2.2. H.264 / H.265 transmission (CGI control)
// 2.2.3. Get UID (User management of video transmission)
// 2.2.4. Response of Get UID
type UID struct {
	// User ID
	UID               string `dcmt:"UID"`

	// Video format
	ImageFormat       string `dcmt:"ImageFormat"`
	// afaik, sensor rating: i.e. 1.3Megapixel
	ImageCaptureMode  string `dcmt:"ImageCaptureMode"`
	// Aspect ratio
	Ratio             string `dcmt:"ratio"`
	Rotation          string `dcmt:"Rotation"`
	// Maximum frame rate
	Maxfps            string `dcmt:"Maxfps"`
	StreamMode        string `dcmt:"StreamMode"`

	// H.265/H.264 bitrate
	IBitrate          string `dcmt:"iBitrate"`
	// H.265/H.264 resolution
	IResolution       string `dcmt:"iResolution"`
	// H.265/H.264 quality
	IQuality          string `dcmt:"iQuality"`

	// Transmission type setting
	SDelivery         string `dcmt:"sDelivery"`

	// Unicast port number
	IUniPort          string `dcmt:"iUniPort"`
	// 1st octet of multicast address
	IMultiAdd1        string `dcmt:"iMultiAdd1"`
	// 2nd octet of multicast address
	IMultiAdd2        string `dcmt:"iMultiAdd2"`
	// 3rd octet of multicast address
	IMultiAdd3        string `dcmt:"iMultiAdd3"`
	// 4th octet of multicast address
	IMultiAdd4        string `dcmt:"iMultiAdd4"`
	// multicast address
	IMultiAdd         string `dcmt:"iMultiAdd"`
	// multicast port number
	IMultiPort        string `dcmt:"iMultiPort"`

	// Audio mode
	AEnable           string `dcmt:"aEnable"`
	// Audio encode
	AEnc              string `dcmt:"aEnc"`
	// Audio bitrate
	ABitrate          string `dcmt:"aBitrate"`
	ABitrate2         string `dcmt:"aBitrate2"`
	ABitrate3         string `dcmt:"aBitrate3"`
	// Audio input interval
	AInterval         string `dcmt:"aInterval"`
	// Audio unicast port number
	AInPort           string `dcmt:"aInPort"`
	// Audio output interval
	AOutInterval      string `dcmt:"aOutInterval"`
	// Audio output port
	AOutPort          string `dcmt:"aOutPort"`
	// Audio output status
	AOutStatus        string `dcmt:"aOutStatus"`
	// Audio output UID
	AOutUID           string `dcmt:"aOutUID"`

	// Event notification port number
	EPort             string `dcmt:"ePort"`
	// Alarm status
	SAlarm            string `dcmt:"sAlarm"`
	// Recording status
	SDrec             string `dcmt:"SDrec"`
	SDrec2            string `dcmt:"SDrec2"`
	// Aux status
	SAUX              string `dcmt:"sAUX"`
	// HTTP port number
	IHttpPort         string `dcmt:"iHttpPort"`

	// Multicast auto for stream (index)
	IMultiAuto_h264   string `dcmt:"iMultiAuto_h264"`
	IMultiAuto_h264_2 string `dcmt:"iMultiAuto_h264_2"`
	IMultiAuto_h264_3 string `dcmt:"iMultiAuto_h264_3"`
	IMultiAuto_h264_4 string `dcmt:"iMultiAuto_h264_4"`

	// Control mode stream (index)
	SRtspMode_h264    string `dcmt:"sRtspMode_h264"`
	SRtspMode_h264_2  string `dcmt:"sRtspMode_h264_2"`
	SRtspMode_h264_3  string `dcmt:"sRtspMode_h264_3"`
	SRtspMode_h264_4  string `dcmt:"sRtspMode_h264_4"`

	// Encode setting for stream (index)
	StreamEncode      string `dcmt:"StreamEncode"`
	StreamEncode_2    string `dcmt:"StreamEncode_2"`
	StreamEncode_3    string `dcmt:"StreamEncode_3"`
	StreamEncode_4    string `dcmt:"StreamEncode_4"`

	// Transmission priority setting for stream (index)
	ITransmit_mode    string `dcmt:"iTransmit_mode"`
	ITransmit_mode_2  string `dcmt:"iTransmit_mode_2"`
	ITransmit_mode_3  string `dcmt:"iTransmit_mode_3"`
	ITransmit_mode_4  string `dcmt:"iTransmit_mode_4"`

	// Smart coding setting for stream (index)
	ISmartCoding      string `dcmt:"iSmartCoding"`
	ISmartCoding_2    string `dcmt:"iSmartCoding_2"`
	ISmartCoding_3    string `dcmt:"iSmartCoding_3"`
	ISmartCoding_4    string `dcmt:"iSmartCoding_4"`
}
