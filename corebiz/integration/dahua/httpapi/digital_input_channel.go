package httpapi

import (
	"time"
	
	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"
)

type FetchDigitalInputChannelListParams struct {}

func (api *APIClient) FetchDigitalInputChannelList(params *FetchDigitalInputChannelListParams) (
	wrapped *response.DigitalInputChannelListWrapper,
	aceI    apicall.APICallEventIntface,
) {
	fnCode := "FetchDigitalInputChannelList"
	fnList := [](func(*FetchDigitalInputChannelListParams)(*response.DigitalInputChannelListWrapper, apicall.APICallEventIntface)){
		api.FetchDigitalInputChannelList__000,
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

func (api *APIClient) FetchDigitalInputChannelList__000(params *FetchDigitalInputChannelListParams) (*response.DigitalInputChannelListWrapper, apicall.APICallEventIntface) {
	_, ev := api.APIV1.GetMachineName()
	
	// inputProxyChannelStatusList, ev := api.APIV1.GetInputProxyChannelStatusList()
	// if inputProxyChannelStatusList == nil {
	// 	return nil, ev
	// }

	// digitalInputChannelList := make([]*response.DigitalInputChannel, 0)
	// if inputProxyChannelStatusList.InputProxyChannelStatusList != nil {
	// 	for _, inputProxyChannelStatus := range inputProxyChannelStatusList.InputProxyChannelStatusList {
	// 		digitalInputChannel := &response.DigitalInputChannel{
	// 			ID                    : inputProxyChannelStatus.ID,
	// 			ProxyProtocol         : "",
	// 			Online                : inputProxyChannelStatus.Online,
	// 			StreamingChannelIDList: make([]string, 0),
	// 			ChannelDetectionResult: inputProxyChannelStatus.ChanDetectResult,
	// 		}

	// 		_sipd := inputProxyChannelStatus.SourceInputPortDescriptor
	// 		if _sipd.AdminProtocol != "" {
	// 			digitalInputChannel.ProxyProtocol = _sipd.AdminProtocol
	// 		} else if _sipd.ProxyProtocol != "" {
	// 			digitalInputChannel.ProxyProtocol = _sipd.ProxyProtocol
	// 		}

	// 		if inputProxyChannelStatus.StreamingProxyChannelIDList != nil {
	// 			for _, streamingProxyChannelId := range inputProxyChannelStatus.StreamingProxyChannelIDList {
	// 				digitalInputChannel.StreamingChannelIDList = append(
	// 					digitalInputChannel.StreamingChannelIDList,
	// 					streamingProxyChannelId.ID,
	// 				)
	// 			}
	// 		}

	// 		digitalInputChannelList = append(digitalInputChannelList, digitalInputChannel)
	// 	}
	// }

	wrapped := &response.DigitalInputChannelListWrapper{
		DigitalInputChannelList: nil,
		// OriginalData          : map[string]any{
		// },
	}

	return wrapped, ev
}
