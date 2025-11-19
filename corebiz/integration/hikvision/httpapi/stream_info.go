package httpapi

import (
	"fmt"
	"time"
	
	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"

	"noname001/corebiz/integration/hikvision/util"
	"noname001/corebiz/integration/hikvision/httpapi/v1/typing"
)

type FetchStreamInfoParams struct {
	ChannelID  string
	StreamType typing.StreamType

	RTSPPort   string
}

func (api *APIClient) FetchStreamInfo(params *FetchStreamInfoParams) (wrapped *response.StreamInfoWrapper, aceI apicall.APICallEventIntface) {
	fnCode := "FetchStreamInfo"
	fnList := [](func(*FetchStreamInfoParams)(*response.StreamInfoWrapper, apicall.APICallEventIntface)){
		api.FetchStreamInfo__000,
	}
	fnCount := len(fnList)

	collector := api.apicallHandler.SpawnCollector(fnCode, fnCount)

	fnIdx, hasMarked := api.apicallHandler.RetrieveMarkedSucceedFunction(fnCode)
	if hasMarked {
		fn := fnList[fnIdx]

		// VVV VVV VVV
		wrapped, aceI = fn(params)
		// ^^^ ^^^ ^^^

		collector.Collect(fnIdx, aceI)
	}

	if !hasMarked {
		for idx, fn := range fnList {

			// VVV VVV VVV
			wrapped, aceI = fn(params)
			// ^^^ ^^^ ^^^

			collector.Collect(idx, aceI)
			if !aceI.IsConsideredError() {
				api.apicallHandler.MarkSucceedFunction(fnCode, idx)
				break
			}

			if (idx + 1) != fnCount {
				// do not spam, wait a bit before trying next function
				time.Sleep(333 * time.Millisecond)
			}
		}
	}

	api.apicallHandler.RetrieveCollector(collector)

	return
}

func (api *APIClient) FetchStreamInfo__000(params *FetchStreamInfoParams) (*response.StreamInfoWrapper, apicall.APICallEventIntface) {
	var (
		wrapped    = &response.StreamInfoWrapper{}
		streamInfo *response.StreamInfo = nil

		evBundle   = apicall.NewBundle("FetchStreamInfo__000")
	)

	// usual value:
	// channelID: 1, 2, 3, ... , 30 31 32
	// streamType: 01, 02, 03
	streamingChannelID := fmt.Sprintf("%s%s", params.ChannelID, params.StreamType)

	streamingChannel, ev := api.APIV1.GetStreamingChannel(streamingChannelID)
	evBundle.AddItem(ev)

	if !ev.IsConsideredError() {
		streamInfo = &response.StreamInfo{}

		streamInfo.ChannelID   = streamingChannel.ID
		streamInfo.ChannelName = streamingChannel.ChannelName
		streamInfo.Enabled     = streamingChannel.Enabled

		streamInfo.VideoEnabled          = streamingChannel.Video.Enabled
		streamInfo.VideoCodecType        = streamingChannel.Video.VideoCodecType
		streamInfo.VideoResolutionWidth  = streamingChannel.Video.VideoResolutionWidth
		streamInfo.VideoResolutionHeight = streamingChannel.Video.VideoResolutionHeight
		streamInfo.VideoFPS              = float32(streamingChannel.Video.MaxFrameRate / 100)

		switch streamingChannel.Video.VideoQualityControlType {
		case "CBR": streamInfo.VideoBitrate = streamingChannel.Video.ConstantBitRate * 1000
		case "VBR": streamInfo.VideoBitrate = streamingChannel.Video.VBRUpperCap * 1000 // TODO: recheck the value
		default:    streamInfo.VideoBitrate = 0
		}

		streamInfo.AudioEnabled   = streamingChannel.Audio.Enabled
		streamInfo.AudioCodecType = streamingChannel.Audio.AudioCompressionType
	}

	wrapped.StreamURL = util.GenerateStreamURL(
		api.hostname, params.RTSPPort,
		api.username, api.password,
		params.ChannelID, string(params.StreamType),
	)
	evBundle.MarkAsPartialSuccess()

	wrapped.StreamInfo = streamInfo

	return wrapped, evBundle
}
