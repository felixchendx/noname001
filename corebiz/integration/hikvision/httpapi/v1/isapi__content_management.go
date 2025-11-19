package v1

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"

	"noname001/corebiz/integration/base/apicall"
)

// 15.2.17 /ISAPI/ContentMgmt/InputProxy/channels/status
func (api *APIClient) GetInputProxyChannelStatusList() (*XML_InputProxyChannelStatusList, *apicall.APICallEvent) {
	ev := apicall.NewEvent("GetInputProxyChannelStatusList")
	defer ev.MarkAsEnded()

	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL + "/ISAPI/ContentMgmt/InputProxy/channels/status",
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


	var successResponse *XML_InputProxyChannelStatusList

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
