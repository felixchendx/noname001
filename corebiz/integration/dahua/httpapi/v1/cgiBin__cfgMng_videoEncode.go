package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"noname001/corebiz/integration/base/apicall"
)

// 4.4VideoEncode

// 4.4.3 GetVideoEncodeConfig
// note: for now, only filters main stream (ExtraFormat[0])
// TODO: test this code against office's NVR
//       at the time of developing this code, dev hardware Dahua NVR is taken for onsite demo
func (api *APIClient) GetAllVideoEncodeConfig() ([]*TXT_VideoEncodeConfig, *apicall.APICallEvent) {
	ev := apicall.NewEvent("GetAllVideoEncodeConfig")
	defer ev.MarkAsEnded()

	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	reqURL := api.baseURL + "/cgi-bin/configManager.cgi?action=getConfig&name=Encode"
	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", reqURL,
		nil,
	)
	if err != nil {
		ev.MarkWithGoError(err)
		return nil, ev
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		ev.MarkWithGoError(err)
		return nil, ev
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ev.MarkWithGoError(err)
		return nil, ev
	}


	var successResponse []*TXT_VideoEncodeConfig

	switch resp.StatusCode {
	case http.StatusOK:
		parsedList := make([]*TXT_VideoEncodeConfig, 0)

		lines := strings.Split(string(body), "\n")
		for _, _line := range lines {
			_line := strings.TrimSpace(_line)
			if _line == "" { continue }

			// splitting "table.Encode[0].ExtraFormat[0].Audio.AudioSource=Coaxial"
			kvParts := strings.Split(_line, "=")
			if len(kvParts) != 2 {
				err := fmt.Errorf("unexpected format: expect 2 kvParts, got '%v'", len(kvParts))
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			// splitting "table.Encode[0].ExtraFormat[0].AudioEnable"
			// splitting "table.Encode[0].ExtraFormat[0].Audio.AudioSource"
			keyParts := strings.Split(kvParts[0], ".")
			switch len(keyParts) {
			case 4, 5 : // pass
			default:
				err := fmt.Errorf("unexpected format: expect 4/5 keyParts, got '%v'", len(keyParts))
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			// discarding non main stream for now
			switch keyParts[2] {
			case "ExtraFormat[0]": // pass
			// see doc for other type
			default:
				continue
			}

			// splitting "Encode[0]"
			encodeParts := strings.Split(keyParts[1], "[")
			if len(encodeParts) != 2 {
				err := fmt.Errorf("unexpected format: expect 2 encodeParts, got '%v'", len(encodeParts))
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			// trimming "0]"
			chanNumber := strings.TrimRight(encodeParts[1], "]")
			if chanNumber == "" {
				err := fmt.Errorf("unexpected format: expect non empty chanNumber, got '%s'", chanNumber)
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			chanNumberInt, convErr := strconv.Atoi(chanNumber)
			if convErr != nil {
				err := fmt.Errorf("unexpected format: non integer chanNumber, got '%s'", chanNumber)
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}


			// TODO: iffy dynamic slice handling
			if len(parsedList) <= chanNumberInt {
				parsedList = append(parsedList, &TXT_VideoEncodeConfig{
					ChannelID: strconv.Itoa(chanNumberInt + 1),
				})
			}

			parsedItem := parsedList[chanNumberInt]

			// non uniform fields: Video.XXX vs VideoEnable
			// presumably the data structure from dahua sth like
			// VideoEncode struct {
			// 		VideoEnable bool
			// 		Video struct {
			// 			resolution xxx
			// 			BitRate    xxx
			// 			...
			// 		}
			// }
			switch keyParts[3] {
			case "VideoEnable": if kvParts[1] == "true" { parsedItem.VideoEnable = true }
			case "AudioEnable": if kvParts[1] == "true" { parsedItem.AudioEnable = true }

			case "Video":
				switch keyParts[4] {
				case "resolution"    : parsedItem.VideoResolution = kvParts[1]
				case "BitRate"       : parsedItem.VideoBitrate, _ = strconv.Atoi(kvParts[1])
				case "BitRateControl": parsedItem.VideoBitrateControl = kvParts[1]
				case "Compression"   : parsedItem.VideoCompression = kvParts[1]
				case "FPS"           : parsedItem.VideoFps, _ = strconv.Atoi(kvParts[1])
				case "GOP"           : parsedItem.VideoGop, _ = strconv.Atoi(kvParts[1])
				case "Height"        : parsedItem.VideoHeight, _ = strconv.Atoi(kvParts[1])
				case "Pack"          : parsedItem.VideoPack = kvParts[1]
				case "Profile"       : parsedItem.VideoProfile = kvParts[1]
				case "Quality"       : parsedItem.VideoQuality, _ = strconv.Atoi(kvParts[1])
				case "QualityRange"  : parsedItem.VideoQualityRange, _ = strconv.Atoi(kvParts[1])
				case "SVCTLayer"     :
				case "Width"         : parsedItem.VideoWidth, _ = strconv.Atoi(kvParts[1])
				default: // warning ?
				}

			case "Audio":
				switch keyParts[4] {
				case "AudioSource" : parsedItem.AudioSource = kvParts[1]
				case "BitRate"     : parsedItem.AudioBitrate, _ = strconv.Atoi(kvParts[1])
				case "Compression" : parsedItem.AudioCompression = kvParts[1]
				case "Depth"       : parsedItem.AudioDepth, _ = strconv.Atoi(kvParts[1])
				case "Frequency"   : parsedItem.AudioFrequency, _ = strconv.Atoi(kvParts[1])
				case "Mode"        : parsedItem.AudioMode, _ = strconv.Atoi(kvParts[1])
				case "Pack"        : parsedItem.AudioPack = kvParts[1]
				case "PacketPeriod": parsedItem.AudioPacketPeriod, _ = strconv.Atoi(kvParts[1])
				default: // warning ?
				}
			}
		}

		successResponse = parsedList

	default:
		failedResponse := &TXT_ResponseStatus{
			RequestURL: reqURL,
			StatusCode: resp.StatusCode,
			StatusMsg : string(body),
		}

		ev.MarkWithAPIError(failedResponse)
		ev.DumpThis(failedResponse.FullError())
		return nil, ev
	}

	return successResponse, ev
}
