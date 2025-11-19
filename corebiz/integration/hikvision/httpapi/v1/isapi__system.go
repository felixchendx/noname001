package v1

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"

	"noname001/corebiz/integration/base/apicall"
)

// 15.10.21 /ISAPI/System/deviceInfo
func (api *APIClient) GetDeviceInfo() (*XML_DeviceInfo, *apicall.APICallEvent) {
	ev := apicall.NewEvent("GetDeviceInfo")

	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL + "/ISAPI/System/deviceinfo",
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


	var successResponse *XML_DeviceInfo

	if resp.StatusCode == http.StatusOK {
		err := xml.Unmarshal(body, &successResponse)
		if err != nil {
			ev.MarkWithGoError(err)
			ev.DumpThis(string(body[:]))
			return nil, ev
		}
	}

	if resp.StatusCode != http.StatusOK {
		var failedResponse *XML_ResponseStatus

		switch resp.StatusCode {
		// TODO: 4xx, 5xx
		default:
			err := xml.Unmarshal(body, &failedResponse)
			if err != nil {
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}
		}

		ev.MarkWithAPIError(failedResponse)
		ev.DumpThis(failedResponse.FullError())
		return nil, ev
	}

	ev.MarkAsEnded()
	return successResponse, ev
}

// 15.10.202 /ISAPI/System/Video/inputs
// analog input channels, should be DVR only endpoint...
func (api *APIClient) GetVideoInputs() (*XML_VideoInput, *apicall.APICallEvent) {
	ev := apicall.NewEvent("GetVideoInputs")
	defer ev.MarkAsEnded()

	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL + "/ISAPI/System/Video/inputs",
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


	var successResponse *XML_VideoInput

	switch resp.StatusCode {
	case http.StatusOK:
		err := xml.Unmarshal(body, &successResponse)
		if err != nil {
			ev.MarkWithGoError(err)
			ev.DumpThis(string(body[:]))
			return nil, ev
		}

	default:
		var failedResponse *XML_ResponseStatus

		err := xml.Unmarshal(body, &failedResponse)
		if err != nil {
			ev.MarkWithGoError(err)
			ev.DumpThis(string(body[:]))
			return nil, ev
		}

		ev.MarkWithAPIError(failedResponse)
		ev.DumpThis(failedResponse.FullError())
		return nil, ev
	}

	return successResponse, ev
}
