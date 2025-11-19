package httpapi

import (
	"time"
	
	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"
)

type FetchAnalogInputChannelListParams struct {}

func (api *APIClient) FetchAnalogInputChannelList(params *FetchAnalogInputChannelListParams) (
	wrapped *response.AnalogInputChannelListWrapper,
	aceI    apicall.APICallEventIntface,
) {
	fnCode := "FetchAnalogInputChannelList"
	fnList := [](func(*FetchAnalogInputChannelListParams)(*response.AnalogInputChannelListWrapper, apicall.APICallEventIntface)){
		api.FetchAnalogInputChannelList__000,
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

func (api *APIClient) FetchAnalogInputChannelList__000(params *FetchAnalogInputChannelListParams) (*response.AnalogInputChannelListWrapper, apicall.APICallEventIntface) {
	var (
		wrapped = &response.AnalogInputChannelListWrapper{}

		evBundle = apicall.NewBundle("FetchAnalogInputChannelList__000")

		len1, len2, longestLen int
	)

	videoEncodeConfigList, ev1 := api.APIV1.GetAllVideoEncodeConfig()
	evBundle.AddItem(ev1)

	time.Sleep(250 * time.Millisecond)

	channelTitleList, ev2 := api.APIV1.GetChannelTitleConfig()
	evBundle.AddItem(ev2)

	if evBundle.IsConsideredError() {
		return nil, evBundle
	}


	if videoEncodeConfigList != nil { len1 = len(videoEncodeConfigList) }
	if channelTitleList != nil { len2 = len(channelTitleList) }

	longestLen = len1
	if longestLen < len2 { longestLen = len2 }

	analogInputChannelList := make([]*response.AnalogInputChannel, 0, longestLen)
	for i := 0; i < longestLen; i++ {
		_analogInputChannel := &response.AnalogInputChannel{}

		if videoEncodeConfigList != nil && i < len1 {
			_item := videoEncodeConfigList[i]

			_analogInputChannel.ID = _item.ChannelID
			// _analogInputChannel.Enabled = _item.VideoEnable || _item.AudioEnable
			_analogInputChannel.Enabled = _item.VideoEnable
			_analogInputChannel.ResDesc = _item.VideoResolution
		}

		if channelTitleList != nil && i < len2 {
			_analogInputChannel.Name = channelTitleList[i]
		}

		analogInputChannelList = append(analogInputChannelList, _analogInputChannel)
	}

	wrapped.AnalogInputChannelList = analogInputChannelList
	wrapped.OriginalData = nil

	return wrapped, evBundle
}
