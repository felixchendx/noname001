package typing

type (
	// Stream type
	VCodec             string
	MultiSensorChannel string

	StreamType string
)

const (
	VCODEC__JPEG   VCodec = "jpeg"
	VCODEC__JPEG_2 VCodec = "jpeg_2"
	VCODEC__JPEG_3 VCodec = "jpeg_3"
	VCODEC__H264   VCodec = "h264"
	VCODEC__H264_2 VCodec = "h264_2"
	VCODEC__H264_3 VCodec = "h264_3"
	VCODEC__H264_4 VCodec = "h264_4"
	VCODEC__H265   VCodec = "h265"
	VCODEC__H265_2 VCodec = "h265_2"
	VCODEC__H265_3 VCodec = "h265_3"
	VCODEC__H265_4 VCodec = "h265_4"

	MSCHAN__CH_1 MultiSensorChannel = "1"
	MSCHAN__CH_2 MultiSensorChannel = "2"
	MSCHAN__CH_3 MultiSensorChannel = "3"
	MSCHAN__CH_4 MultiSensorChannel = "4"

	STREAM_1 StreamType = "1"
	STREAM_2 StreamType = "2"
	STREAM_3 StreamType = "3"
	STREAM_4 StreamType = "4"
)
