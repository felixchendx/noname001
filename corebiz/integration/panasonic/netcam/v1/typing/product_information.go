package typing

// 9.1. Get product information
type ProductInformation struct {
	// MAC address
	Mac                string `dcmt:"MAC"`
	// Serial number
	Serial             string `dcmt:"SERIAL"`
	// Software version
	Version            string `dcmt:"VERSION"`
	// Product name e.g.)WV-S2531L
	Name               string `dcmt:"NAME"`

	// Status of SD card recording stream 1
	SDrec              string `dcmt:"SDrec"`
	// Status of SD card recording stream 2
	SDrec2             string `dcmt:"SDrec2"`

	// Alarm status (CH1)
	SAlarm             string `dcmt:"sAlarm"`
	// AUX output status
	SAUX               string `dcmt:"sAUX"`

	// --- not in doc section's doc, see uid section ---
	EPort              string `dcmt:"ePort"`

	// Current value of audio mode setting
	AEnable            string `dcmt:"aEnable"`
	// Audio encoder setup
	AEnc               string `dcmt:"aEnc"`
	// Current value of audio bit rate setting (G.726)
	ABitrate           string `dcmt:"aBitrate"`
	// Current value of audio bit rate setting (G.711)
	ABitrate2          string `dcmt:"aBitrate2"`
	// Current value of audio bit rate setting (AAC-LC)
	ABitrate3          string `dcmt:"aBitrate3"`
	// Current value of audio input interval setting (from camera to PC)
	AInInterval        string `dcmt:"aInInterval"`
	// Current value of audio output interval setting (from PC to camera)
	AOutInterval       string `dcmt:"aOutInterval"`
	// Current value of audio output port setting (from PC to camera)
	AOutPort           string `dcmt:"aOutPort"`
	// Status of audio output function
	AOutStatus         string `dcmt:"aOutStatus"`
	// UID that is transmitting "audio output"
	AOutUID            string `dcmt:"aOutUID"`

	// --- not in this section's doc, see uid section ---
	AInPort_h264       string `dcmt:"aInPort_h264"`
	AInPort_h264_2     string `dcmt:"aInPort_h264_2"`
	AInPort_h264_3     string `dcmt:"aInPort_h264_3"`
	AInPort_h264_4     string `dcmt:"aInPort_h264_4"`

	// --- not in this section's doc, see uid section ---
	SRtspMode_h264     string `dcmt:"sRtspMode_h264"`
	SRtspMode_h264_2   string `dcmt:"sRtspMode_h264_2"`
	SRtspMode_h264_3   string `dcmt:"sRtspMode_h264_3"`
	SRtspMode_h264_4   string `dcmt:"sRtspMode_h264_4"`

	// Current value of the Image capture mode setting
	ImageCaptureMode   string `dcmt:"ImageCaptureMode"`

	// Aspect ratio
	Ratio              string `dcmt:"ratio"`
	// --- ??? ---
	Ratio_sub          string `dcmt:"ratio_sub"`

	// --- no desc, see uid section ---
	Maxfps             string `dcmt:"Maxfps"`
	// 1(fixed)
	StreamMode         string `dcmt:"StreamMode"`
	// --- no desc, see uid section ---
	Rotation           string `dcmt:"Rotation"`

	// Encode setting of stream(idx)
	StreamEncode       string `dcmt:"StreamEncode"`
	StreamEncode_2     string `dcmt:"StreamEncode_2"`
	StreamEncode_3     string `dcmt:"StreamEncode_3"`
	StreamEncode_4     string `dcmt:"StreamEncode_4"`

	// Current value of the stream(1) setting
	// 0: OFF, 1: ON
	ITransmit_h264     string `dcmt:"iTransmit_h264"`
	// Current value of stream(1) setting.
	// uni : Unicast (Auto)
	// multi : Multicast
	// uni_manual : Unicast (Manual)
	SDelivery_h264     string `dcmt:"sDelivery_h264"`
	// Current value of the Stream(1) bandwidth
	// setting.
	// 64 : 64kbps,
	// 128 : 128 kbps,
	// 256 : 256 kbps,
	// 384: 384 kbps,
	// 512: 512 kbps,
	// 768: 768 kbps,
	// 1024: 1024 kbps,
	// 1536: 1536 kbps,
	// 2048: 2048 kbps,
	// 3072: 3072 kbps,
	// 4096: 4096 kbps,
	// 6144: 6144 kbps
	// 8192: 8192 kbps,
	// 10240: 10240 kbps,
	// 12288: 12288 kbps
	// 14336: 14336 kbps
	// 16384: 16384 kbps
	// 20480: 20480 kbps
	// 24576: 24576 kbps
	IBitrate_h264      string `dcmt:"iBitrate_h264"`
	// Current value of the Stream(1) resolution
	// setting.
	// Resolution to be set (4:3)
	// 320 : QVGA
	// 400 : 400x300
	// 640 : VGA
	// 1280 : 1280 x 960
	// 2048 : 2048 x 1536
	// 800 800 x 600
	// 1600: 1600x1200
	// 2560: 2560x1920
	// 3072: 3072x2304
	// Resolution to be set (16:9)
	// 640 : 640 x 360
	// 320 : 320 x 180
	// 1280 : 1280 x 720
	// 1920 : 1920 x 1080
	// 2560: 2560x1440
	// 3072: 3072x1728
	// 3840: 3840x2160
	// Resolution to set (1:1)
	// 640 : 640 x 640
	// 320 : 320 x 320
	// 1280 : 1280 x 1280
	// 2192 : 2192 x 2192
	// 2992 : 2992 x 2992
	IResolution_h264   string `dcmt:"iResolution_h264"`
	// Current value of the Stream(1) quality setting
	// fine : Fine
	// normal: Normal
	// low : Low
	// 0, 1, 2, 3, 4, 5, 6, 7, 8, 9：10 step setting when
	// VBR
	IQuality_h264      string `dcmt:"iQuality_h264"`
	// Multicast auto start stream(1)
	IMultiAuto_h264    string `dcmt:"iMultiAuto_h264"`
	// Stream(1) priority setting
	ITransmit_mode     string `dcmt:"iTransmit_mode"`
	// Stream(1) Smartcoding setting
	ISmartCoding       string `dcmt:"iSmartCoding"`

	// --- same as above, but Stream(2)
	ITransmit_h264_2   string `dcmt:"iTransmit_h264_2"`
	SDelivery_h264_2   string `dcmt:"sDelivery_h264_2"`
	IBitrate_h264_2    string `dcmt:"iBitrate_h264_2"`
	IResolution_h264_2 string `dcmt:"iResolution_h264_2"`
	IQuality_h264_2    string `dcmt:"iQuality_h264_2"`
	IMultiAuto_h264_2  string `dcmt:"iMultiAuto_h264_2"`
	ITransmit_mode_2   string `dcmt:"iTransmit_mode_2"`
	ISmartCoding_2     string `dcmt:"iSmartCoding_2"`

	// --- same as above, but Stream(3)
	ITransmit_h264_3   string `dcmt:"iTransmit_h264_3"`
	SDelivery_h264_3   string `dcmt:"sDelivery_h264_3"`
	IBitrate_h264_3    string `dcmt:"iBitrate_h264_3"`
	IResolution_h264_3 string `dcmt:"iResolution_h264_3"`
	IQuality_h264_3    string `dcmt:"iQuality_h264_3"`
	IMultiAuto_h264_3  string `dcmt:"iMultiAuto_h264_3"`
	ITransmit_mode_3   string `dcmt:"iTransmit_mode_3"`
	ISmartCoding_3     string `dcmt:"iSmartCoding_3"`

	// --- same as above, but Stream(4)
	ITransmit_h264_4   string `dcmt:"iTransmit_h264_4"`
	SDelivery_h264_4   string `dcmt:"sDelivery_h264_4"`
	IBitrate_h264_4    string `dcmt:"iBitrate_h264_4"`
	IResolution_h264_4 string `dcmt:"iResolution_h264_4"`
	IQuality_h264_4    string `dcmt:"iQuality_h264_4"`
	IMultiAuto_h264_4  string `dcmt:"iMultiAuto_h264_4"`
	ITransmit_mode_4   string `dcmt:"iTransmit_mode_4"`
	ISmartCoding_4     string `dcmt:"iSmartCoding_4"`

	// --- ??? ---
	PRODUCT_NAME       string `dcmt:"PRODUCT_NAME"`
	// --- ??? ---
	DESTINATION        string `dcmt:"DESTINATION"`


	// --- --- ---
	// these fields below are in doc... but...

	// Firmware name e.g)s7130
	Firmware           string `dcmt:"FIRMWARE"`

	// Alarm status (CH2)
	SAlarm2            string `dcmt:"sAlarm2"`
	// Alarm status (CH3)
	SAlarm3            string `dcmt:"sAlarm3"`
	// Alarm status (CH4)
	SAlarm4            string `dcmt:"sAlarm4"`
}
