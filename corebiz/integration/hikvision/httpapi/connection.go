package httpapi

import (
	"time"
	
	"noname001/corebiz/integration/base/apicall"
)

func (api *APIClient) TestConnection() (connOK bool, aceI apicall.APICallEventIntface) {
	fnCode := "TestConnection"
	fnList := [](func()(bool, apicall.APICallEventIntface)){
		api.TestConnection__000,
	}
	fnCount := len(fnList)

	collector := api.apicallHandler.SpawnCollector(fnCode, fnCount)

	fnIdx, hasMarked := api.apicallHandler.RetrieveMarkedSucceedFunction(fnCode)
	if hasMarked {
		fn := fnList[fnIdx]

		// VVV VVV VVV
		connOK, aceI = fn()
		// ^^^ ^^^ ^^^

		collector.Collect(fnIdx, aceI)
	}

	if !hasMarked {
		for idx, fn := range fnList {

			// VVV VVV VVV
			connOK, aceI = fn()
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

func (api *APIClient) TestConnection__000() (bool, apicall.APICallEventIntface) {
	_, aceI := api.APIV1.GetDeviceInfo()
	connOK := !aceI.IsConsideredError()
	return connOK, aceI
}
