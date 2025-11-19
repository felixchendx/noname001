package httpapi

import (
	"time"

	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"

	"noname001/corebiz/integration/dahua/util"
	// "noname001/corebiz/integration/dahua/httpapi/v1/typing"
)

type FetchStreamInfoParams struct {
	ChannelID  string
	// StreamType typing.StreamType

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

// REVERTED FROM BEFORE jun-008 + adjusted to new struct, see bottom for context
func (api *APIClient) FetchStreamInfo__000(params *FetchStreamInfoParams) (*response.StreamInfoWrapper, apicall.APICallEventIntface) {
	// streamInfo, ev := api.APIV1.GetStreamInfo()
	// if streamInfo == nil {
	// 	return nil, ev
	// }
	wrapped := &response.StreamInfoWrapper{
		StreamInfo: &response.StreamInfo{},
		OriginalData: make(map[string]any),
	}
	evBundle := apicall.NewBundle("FetchStreamInfo__000")
	evBundle.MarkAsPartialSuccess()

	streamURL := util.GenerateStreamURL(
		api.hostname, params.RTSPPort,
		api.username, api.password,
		params.ChannelID, "0",
	)

	videoConfig, ev1:=  api.APIV1.GetVideoEncodeConfig(params.ChannelID)
	evBundle.AddItem(ev1)
	if (!ev1.IsConsideredError()) && (videoConfig != nil){
		wrapped.OriginalData["FetchStreamInfo__000"] = videoConfig

		wrapped.StreamInfo.VideoEnabled = videoConfig.VideoEnable
		wrapped.StreamInfo.VideoCodecType = videoConfig.VideoCompression
		wrapped.StreamInfo.VideoResolutionWidth = videoConfig.VideoWidth
		wrapped.StreamInfo.VideoResolutionHeight = videoConfig.VideoHeight
		wrapped.StreamInfo.VideoBitrate = videoConfig.VideoBitrate * 1000 
		wrapped.StreamInfo.VideoFPS = float32(videoConfig.VideoFps)

		wrapped.StreamInfo.AudioEnabled = videoConfig.AudioEnable
		wrapped.StreamInfo.AudioCodecType = videoConfig.AudioCompression


	}

	wrapped.StreamURL = streamURL
	wrapped.StreamInfo.ChannelID = params.ChannelID

	return wrapped, evBundle
}

// TODO: case-yogya, crash loop NPE on videoConfig nil
// only ev is not inaf, the videoConfig can still be nil
// might need to rework the response parsing
// TODO: extra stream support after rework response parsing
// func (api *APIClient) FetchStreamInfo__000(params *FetchStreamInfoParams) (*response.StreamInfoWrapper, apicall.APICallEventIntface) {
// 	var (
// 		wrapped    = &response.StreamInfoWrapper{}
// 		streamInfo *response.StreamInfo = nil

// 		evBundle   = apicall.NewBundle("FetchStreamInfo__000")
// 	)

// 	videoConfig, ev1 := api.APIV1.GetVideoEncodeConfig(params.ChannelID)
// 	evBundle.AddItem(ev1)

// 	if !ev1.IsConsideredError() {
// 		streamInfo = &response.StreamInfo{}

// 		streamInfo.VideoEnabled          = videoConfig.VideoEnable
// 		streamInfo.VideoCodecType        = videoConfig.VideoCompression
// 		streamInfo.VideoResolutionWidth  = videoConfig.VideoWidth
// 		streamInfo.VideoResolutionHeight = videoConfig.VideoHeight
// 		streamInfo.VideoBitrate          = videoConfig.VideoBitrate * 1000 
// 		streamInfo.VideoFPS              = float32(videoConfig.VideoFps)

// 		streamInfo.AudioEnabled   = videoConfig.AudioEnable
// 		streamInfo.AudioCodecType = videoConfig.AudioCompression
// 	}

// 	wrapped.StreamURL = util.GenerateMainStreamURL(
// 		api.hostname, params.RTSPPort,
// 		api.username, api.password,
// 		params.ChannelID,
// 	)
// 	evBundle.MarkAsPartialSuccess()

// 	wrapped.StreamInfo = streamInfo
// 	// wrapped.OriginalData["FetchStreamInfo__000"] = videoConfig

// 	return wrapped, evBundle
// }
