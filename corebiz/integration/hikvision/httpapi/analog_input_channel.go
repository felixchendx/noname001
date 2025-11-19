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
	videoInput, ev := api.APIV1.GetVideoInputs()
	if ev.IsConsideredError() {
		return nil, ev
	}

	analogInputChannelList := make([]*response.AnalogInputChannel, 0)
	if videoInput.VideoInputChannelList != nil {
		for _, videoInputChannel := range videoInput.VideoInputChannelList {
			analogInputChannelList = append(analogInputChannelList, &response.AnalogInputChannel{
				ID         : videoInputChannel.ID,
				Enabled    : videoInputChannel.VideoInputEnabled,
				Name       : videoInputChannel.Name,
				VideoFormat: videoInputChannel.VideoFormat,
				PortType   : videoInputChannel.PortType,
				ResDesc    : videoInputChannel.ResDesc,
			})
		}
	}

	wrapped := &response.AnalogInputChannelListWrapper{
		AnalogInputChannelList: analogInputChannelList,
		// OriginalData          : map[string]any{
		// },
	}

	return wrapped, ev
}
