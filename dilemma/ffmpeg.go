package dilemma

import (
	"strconv"
)

const (
	FFMPEG_BIN_PATH = "/usr/bin/ffmpeg"
)

// type VideoEncoder string
// type VideoDecoder string

// type AudioEncoder string
// type AudioDecoder string

type VideoCodec string
type AudioCodec string

type VideoBitrate int // in bit/s
type AudioBitrate int // in bit/s

type CompressionRatio int // 0-100
type BitrateUnit string

const (
	// VIDEO_ENCODER_H264 VideoEncoder = "libx264"
	// VIDEO_ENCODER_H265 VideoEncoder = "libx265"
	
	// VIDEO_DECODER_H264 VideoDecoder = "H264"
	// VIDEO_DECODER_H265 VideoDecoder = "HEVC"

	VIDEO_CODEC_H264 VideoCodec = "h264" // TODO: x264
	VIDEO_CODEC_H265 VideoCodec = "h265"

	AUDIO_CODEC_OPUS AudioCodec = "opus"
	AUDIO_CODEC_AAC  AudioCodec = "aac"

	BITRATE_UNIT_kbps BitrateUnit = "kbps"
)

func IntToVideoBitrate(bitrate int) (vb VideoBitrate) {
	// TODO : validate
	return VideoBitrate(bitrate)
}

func StringToVideoCodec(codec string) (typedCodec VideoCodec) {
	switch codec {
		case string(VIDEO_CODEC_H264): typedCodec = VIDEO_CODEC_H264
		case string(VIDEO_CODEC_H265): typedCodec = VIDEO_CODEC_H265
		default: typedCodec = VIDEO_CODEC_H264
	}

	return
}

func StringToAudioCodec(codec string) (typedCodec AudioCodec) {
	switch codec {
		case string(AUDIO_CODEC_AAC): typedCodec = AUDIO_CODEC_AAC
		case string(AUDIO_CODEC_OPUS): typedCodec = AUDIO_CODEC_OPUS
		default: typedCodec = AUDIO_CODEC_AAC
	}

	return
}

func (bitrate VideoBitrate) String() (string) {
	return strconv.Itoa(int(bitrate))
}

func (bitrate AudioBitrate) String() (string) {
	return strconv.Itoa(int(bitrate))
}

func (ratio CompressionRatio) String() (string) {
	return strconv.Itoa(int(ratio))
}
