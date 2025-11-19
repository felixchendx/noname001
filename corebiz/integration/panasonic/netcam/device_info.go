package netcam

import (
	"time"

	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"
)

type FetchDeviceInfoParams struct {}

func (api *APIClient) FetchDeviceInfo(params *FetchDeviceInfoParams) (wrapped *response.DeviceInfoWrapper, aceI apicall.APICallEventIntface) {
	fnCode := "FetchDeviceInfo"
	fnList := [](func(*FetchDeviceInfoParams)(*response.DeviceInfoWrapper, apicall.APICallEventIntface)){
		api.FetchDeviceInfo__000,
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

func (api *APIClient) FetchDeviceInfo__000(params *FetchDeviceInfoParams) (*response.DeviceInfoWrapper, apicall.APICallEventIntface) {
	productInfo, ev := api.APIV1.GetProductInformation()
	if ev.IsConsideredError() {
		return nil, ev
	}

	wrapped := &response.DeviceInfoWrapper{}
	wrapped.DeviceInfo = &response.DeviceInfo{
		DeviceName: productInfo.Name,
		DeviceID  : "",
		Model     : productInfo.PRODUCT_NAME,
		DeviceType: "Network Camera",
	}

	if wrapped.DeviceInfo.Model == "" {
		wrapped.DeviceInfo.Model = productInfo.Name
	}

	wrapped.OriginalData = make(map[string]any)
	wrapped.OriginalData["productInfo"] = productInfo

	return wrapped, ev
}
