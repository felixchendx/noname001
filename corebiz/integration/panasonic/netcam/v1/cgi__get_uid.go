package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"noname001/corebiz/integration/base/apicall"

	"noname001/corebiz/integration/panasonic/netcam/v1/typing"
)

// 2.2.3. Get UID (User management of video transmission)
// [URL] /cgi-bin/getuid?FILE=2&vcodec=< Value>&reply=info[&ch=<Value>]
// [Method] GET
// [Access level] 3
func (api *APIClient) GetUID(vcodec typing.VCodec, ch *typing.MultiSensorChannel) (*typing.UID, *apicall.APICallEvent) {
	acev := apicall.NewEvent("GetUID")
	defer acev.MarkAsEnded()

	reqContext, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	reqURI := fmt.Sprintf("/cgi-bin/getuid?FILE=2&vcodec=%s&reply=info", vcodec)
	if ch != nil {
		reqURI += fmt.Sprintf("&ch=%s", ch)
	}

	req, err := http.NewRequestWithContext(
		reqContext,
		"GET", api.baseURL + reqURI,
		nil,
	)
	if err != nil {
		acev.MarkWithGoError(err)
		return nil, acev
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		acev.MarkWithGoError(err)
		return nil, acev
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		acev.MarkWithGoError(err)
		return nil, acev
	}


	var successResponse = &typing.UID{}

	switch resp.StatusCode {
	case http.StatusOK:
		decodingReport := api.decodeSuccessResponse(string(body), successResponse)
		if decodingReport.err != nil {
			acev.MarkWithGoError(decodingReport.err)
			acev.DumpThis(string(body[:]))
			return nil, acev
		}

	case http.StatusUnauthorized:
		fallthrough
	case http.StatusForbidden:
		fallthrough
	case http.StatusNotFound:
		var failedResponse = &typing.FailedResponse{}

		// temp
		// api.decodeFailedResponse(string(body), failedResponse)
		failedResponse.Status = string(resp.StatusCode)
		failedResponse.Message = string(body)

		acev.MarkWithAPIError(failedResponse)
		acev.DumpThis(failedResponse.FullError())
		return nil, acev

	default:
		var unknownResponse = &typing.UnknownResponse{
			Status : fmt.Sprintf("Unknown - %s", resp.StatusCode),
			Message: "Unknown response",
		}

		acev.MarkWithAPIError(unknownResponse)
		acev.DumpThis(unknownResponse.FullError())
		acev.DumpThis(string(body[:]))
		return nil, acev
	}

	return successResponse, acev
}
