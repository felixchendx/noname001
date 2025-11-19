package httpapi

import (
	"time"

	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"
	// v1 "noname001/corebiz/integration/dahua/httpapi/v1"
)

type FetchDeviceInfoParams struct{}

func (api *APIClient) FetchDeviceInfo(params *FetchDeviceInfoParams) (wrapped *response.DeviceInfoWrapper, aceI apicall.APICallEventIntface) {
	fnCode := "FetchDeviceInfo"
	fnList := [](func(*FetchDeviceInfoParams) (*response.DeviceInfoWrapper, apicall.APICallEventIntface)){
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
	eventBundle := apicall.NewBundle("FetchDeviceInfo__000")
	var wrapped *response.DeviceInfoWrapper = &response.DeviceInfoWrapper{
		DeviceInfo:   &response.DeviceInfo{},
		// OriginalData: make(map[string]any),
	}
	// var deviceInfoBundle *v1.TXT_DeviceInfo = &v1.TXT_DeviceInfo{}

	deviceName, event1 := api.APIV1.GetMachineName()
	eventBundle.AddItem(event1)
	if !event1.IsConsideredError() {
		wrapped.DeviceInfo.DeviceName = deviceName
		// deviceInfoBundle.DeviceName = deviceName

		eventBundle.MarkAsPartialSuccess()
	}

	time.Sleep(250 * time.Millisecond)
	deviceType, event2 := api.APIV1.GetDeviceType()
	eventBundle.AddItem(event2)
	if !event2.IsConsideredError() {
		wrapped.DeviceInfo.Model = deviceType
		// deviceInfoBundle.DeviceType = deviceType

		eventBundle.MarkAsPartialSuccess()
	}

	// time.Sleep(250 * time.Millisecond)
	// deviceManufacture, event4 := api.APIV1.GetVendor()
	// eventBundle.AddItem(event4)
	// if (deviceManufacture != "") || (!event4.IsConsideredError()) {
	// 	deviceInfoBundle.Manufacturer = deviceManufacture
	// }

	// time.Sleep(250 * time.Millisecond)
	// hardwareVersion, event5 := api.APIV1.GetHardwareVersion()
	// eventBundle.AddItem(event5)
	// if (hardwareVersion != "") || (!event5.IsConsideredError()) {
	// 	deviceInfoBundle.HardwareVersion = hardwareVersion
	// }

	// time.Sleep(250 * time.Millisecond)
	// softwareVersion, event6 := api.APIV1.GetSoftwareVersion()
	// eventBundle.AddItem(event6)
	// if (softwareVersion != nil) || (!event6.IsConsideredError()) {
	// 	deviceInfoBundle.SoftwareVersion = softwareVersion
	// }

	// time.Sleep(250 * time.Millisecond)
	// systemInfo, event7 := api.APIV1.GetSystemInfo()
	// eventBundle.AddItem(event7)
	// if (systemInfo != nil) || (!event7.IsConsideredError()) {
	// 	deviceInfoBundle.SystemInfo = systemInfo
	// }

	// wrapped.OriginalData["FetchDeviceInfo__000"] = deviceInfoBundle

	if eventBundle.IsConsideredError() {
		return nil, eventBundle
	}

	return wrapped, eventBundle
}
