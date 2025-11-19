package v1

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	recordTypeReguler string = "reguler_record"
	recordTypeMotion  string = "motion_detection_sensor"
	recordTypeAlarm   string = "alarm_record"
	headMain          string = "MainFormat"
	MaxAmountChannels int    = 16
)

// Parse value extracts a single value from the HTTP response body Dahua.
//
// Parameters:
//   - body: The body of the HTTP response (example = "version=V1.0").
//   - key : The key of the value to retrieve (example = "version").
//
// Returns:
//   - output: The value associated with the key (example = "V1.0").
//   - err   : nil if successful, or an error if it fails.
func (api *APIClient) parseSingleValueString(body string, key string) (output string, err error) {
	result := strings.Split(body, "=")
	if (len(result) == 2) && (result[0] == key) && (len(result[1]) > 0) {
		output = strings.TrimSpace(result[1])
		return output, nil
	}

	return "", fmt.Errorf("failed parse single value in key filter %s", key)
}

// Parse software version value from the HTTP response body Dahua.
//
// Parameters:
//   - body: The body of the HTTP response (example = "version=4.001.0000000.17.R,build:2021-12-30 21:09:34").
//
// Returns:
//   - output: The value associated with the key (example = ["4.001.0000000.17.R", "2021-12-30 21:09:34"]).
//   - err   : nil if successful, or an error if it fails.
func (api *APIClient) parseSoftwareVersion(body string) (version, releaseDate string, err error) {
	parts := strings.Split(body, ",")
	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			keyValue = strings.SplitN(part, ":", 2)
		}
		if len(keyValue) != 2 {
			return "", "", fmt.Errorf("unexpected software version format v1 response body. body = %s", body)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		switch key {
		case "version":
			version = value
		case "build":
			releaseDate = value
		default:
			return "", "", fmt.Errorf("unexpected software version format v1 response body. body = %s", body)
		}
	}

	return version, releaseDate, nil
}

// Parse system info Value from the HTTP response body Dahua.
//
// Parameters:
//   - body: A string representing the body of the HTTP response. Example:
//     "processor=ST7108
//     serialNumber=8E042A0PAZ91103
//     updateSerial=DH-XVR5216AN-I3"
//
// Returns:
//   - output: The value associated with the key (example = ["4.001.0000000.17.R", "2021-12-30 21:09:34"]).
//   - err   : nil if successful, or an error if it fails.
func (api *APIClient) parseSystemInfo(statement string) (processor, serialNumber, updateSerial string, err error) {
	lines := strings.Split(statement, "\n")

	processor = ""
	serialNumber = ""
	updateSerial = ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("unexpected system info format v1 response body. body = %s", statement)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "processor":
			processor = value
		case "serialNumber":
			serialNumber = value
		case "updateSerial":
			updateSerial = value
		default:
			return "", "", "", fmt.Errorf("unexpected system info format v1 response body. body = %s", statement)
		}
	}

	return processor, serialNumber, updateSerial, nil
}

// Parse Video Encode Config from the HTTP response body Dahua.
//
// Parameters:
//   - channelID: A string channel number. Example: "4"
//   - body: A string representing the body of the HTTP response. Example:
//     "table.Encode[4].ExtraFormat[0].Audio.AudioSource=Coaxial
//     table.Encode[4].ExtraFormat[0].Audio.BitRate=10
//     table.Encode[4].ExtraFormat[0].Audio.Compression=G.711A
//     ... "
//
// Returns:
//   - output: The value associated with the key (example = ["", "4", ....]).
//   - err   : nil if successful, or an error if it fails.
func (api *APIClient) parseVideoEncodeConfig(channelID, body string) (*TXT_VideoEncodeConfig, error) {
	var VideoEncodeConfig *TXT_VideoEncodeConfig = &TXT_VideoEncodeConfig{}
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ".")
		if len(parts) < 3 {
			return nil, fmt.Errorf("unexpected format video config encode format v1 response body. line = %s", line)
		}

		recordType, err := api.parserRegulerRecord(parts[2])
		if err != nil {
			return nil, err
		} else if recordType == "" {
			continue
		}

		VideoEncodeConfig.ChannelID = channelID
		VideoEncodeConfig.ChannelRecordType = recordType
		err = api.parserVideoAudioChannelsInfo(VideoEncodeConfig, line, channelID)
		if err != nil {
			return nil, err
		}
	}

	return VideoEncodeConfig, nil
}

// Parse record type (must regular_record) camera form Api GetVideoEncodeConfig
//
// Parameters:
//   - statement: A string input after split. Example:
//     "MainFormat[0].Audio.BitRate=10"
//
// Returns:
//   - output: The value reguler record (example = 0).
//   - err   : nil if recordType == 0, or an error if recordType != 0.
func (api *APIClient) parserRegulerRecord(statement string) (recordType string, err error) {
	parts := strings.Split(statement, "[")
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected record type of video config encode format v1 response body. data = %s", statement)
	}

	if parts[0] != headMain {
		return "", nil
	}

	numberPart := strings.TrimSuffix(parts[1], "]")
	number, err := strconv.Atoi(numberPart)
	if err != nil {
		return "", fmt.Errorf("failed convert format number record type str to int. err = %v; data = %s", err, statement)

	}

	if number != 0 {
		return "", nil
	}

	return recordTypeReguler, nil
}

// Parse string value from Api GetVideoEncodeConfig about audio/video value based filter
//
// Parameters:
//   - statement: A string input. Example: "Audio.Compression=G.711A"
//   - filter   : A key of statement. Example: "Audio.Compression"
//   - output   : A string output. Example: "G.711A"
func (api *APIClient) parseAudioOrVideoValueString(statement string, filter string, output *string) {
	result := strings.Split(statement, filter)
	if (len(result) == 2) && (len(result[1]) > 0) {
		*output = result[1]
	}
}

// Parse int value from Api GetVideoEncodeConfig about audio/video value based filter
//
// Parameters:
//   - statement: A string input. Example: "Audio.BitRate=10"
//   - filter   : A key of statement. Example: "Audio.BitRate"
//   - output   : A string output. Example: 10
func (api *APIClient) parseAudioOrVideoValueInt(statement string, filter string, output *int) {
	result := strings.Split(statement, filter)
	if (len(result) == 2) && (len(result[1]) > 0) {
		*output, _ = strconv.Atoi(result[1])
	}
}

// Parse bool value from Api GetVideoEncodeConfig about audio/video value based filter
//
// Parameters:
//   - statement: A string input. Example: ".VideoEnable=true"
//   - filter   : A key of statement. Example: "VideoEnable"
//   - output   : A string output. Example: true/false
func (api *APIClient) parseAudioOrVideoValueBool(statement string, filter string, output *bool) {
	result := strings.Split(statement, filter)
	if (len(result) == 2) && (len(result[1]) > 0) {
		switch result[1] {
		case "true":
			*output = true
		case "false":
			*output = false
		}
	}
}

// Parse all about video/audio config value from Api GetVideoEncodeConfig
//
// Parameters:
//   - channelInfo : A struct output info channel.
//   - statement   : A string input. Example: "table.Encode[4].MainFormat[1].VideoEnable=true"
//   - channelId   : A string channel id. Example: "4"
//
// Return:
//   - err   : nil if successful, or an error if it fails.
func (api *APIClient) parserVideoAudioChannelsInfo(channelInfo *TXT_VideoEncodeConfig, statement string, channelId string) error {
	key := fmt.Sprintf("table.Encode[%s].%s[0].", channelId, headMain)
	filter := strings.Split(statement, key)
	if len(filter) != 2 {
		return fmt.Errorf("unexpected line video config encode format v1 response body. data = %s", statement)
	}

	api.parseAudioOrVideoValueString(filter[1], "Audio.AudioSource=", &channelInfo.AudioSource)
	api.parseAudioOrVideoValueInt(filter[1], "Audio.BitRate=", &channelInfo.AudioBitrate)
	api.parseAudioOrVideoValueString(filter[1], "Audio.Compression=", &channelInfo.AudioCompression)
	api.parseAudioOrVideoValueInt(filter[1], "Audio.Depth=", &channelInfo.AudioDepth)
	api.parseAudioOrVideoValueInt(filter[1], "Audio.Frequency=", &channelInfo.AudioFrequency)
	api.parseAudioOrVideoValueInt(filter[1], "Audio.Mode=", &channelInfo.AudioMode)
	api.parseAudioOrVideoValueString(filter[1], "Audio.Pack=", &channelInfo.AudioPack)
	api.parseAudioOrVideoValueInt(filter[1], "Audio.PacketPeriod=", &channelInfo.AudioPacketPeriod)
	api.parseAudioOrVideoValueBool(filter[1], "AudioEnable=", &channelInfo.AudioEnable)
	api.parseAudioOrVideoValueString(filter[1], "Video.resolution=", &channelInfo.VideoResolution)
	api.parseAudioOrVideoValueInt(filter[1], "Video.BitRate=", &channelInfo.VideoBitrate)
	api.parseAudioOrVideoValueString(filter[1], "Video.BitRateControl=", &channelInfo.VideoBitrateControl)
	api.parseAudioOrVideoValueString(filter[1], "Video.Compression=", &channelInfo.VideoCompression)
	api.parseAudioOrVideoValueInt(filter[1], "Video.FPS=", &channelInfo.VideoFps)
	api.parseAudioOrVideoValueInt(filter[1], "Video.GOP=", &channelInfo.VideoGop)
	api.parseAudioOrVideoValueInt(filter[1], "Video.Height=", &channelInfo.VideoHeight)
	api.parseAudioOrVideoValueInt(filter[1], "Video.Width=", &channelInfo.VideoWidth)
	api.parseAudioOrVideoValueString(filter[1], "Video.Pack=", &channelInfo.VideoPack)
	api.parseAudioOrVideoValueString(filter[1], "Video.Profile=", &channelInfo.VideoProfile)
	api.parseAudioOrVideoValueInt(filter[1], "Video.Quality=", &channelInfo.VideoQuality)
	api.parseAudioOrVideoValueInt(filter[1], "Video.QualityRange=", &channelInfo.VideoQualityRange)
	api.parseAudioOrVideoValueBool(filter[1], "VideoEnable=", &channelInfo.VideoEnable)

	return nil
}
