package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"noname001/corebiz/integration/base/apicall"
)

// 4.7ChannelTitle

// 4.7.1 GetChannelTitleConfig
// TODO: test this code against office's NVR
//       at the time of developing this code, dev hardware Dahua NVR is taken for onsite demo
func (api *APIClient) GetChannelTitleConfig() ([]string, *apicall.APICallEvent) {
	ev := apicall.NewEvent("GetChannelTitleConfig")
	defer ev.MarkAsEnded()

	reqCtx, reqC := context.WithTimeout(api.context, api.httpTimeout)
	defer reqC()

	reqURL := api.baseURL + "/cgi-bin/configManager.cgi?action=getConfig&name=ChannelTitle"
	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", reqURL,
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


	var successResponse []string

	switch resp.StatusCode {
	case http.StatusOK:
		parsedList := make([]string, 0)

		lines := strings.Split(string(body), "\n")
		for _, _line := range lines {
			_line := strings.TrimSpace(_line)
			if _line == "" { continue }

			// splitting "table.ChannelTitle[0].Name=Channel1"
			kvParts := strings.Split(_line, "=")
			if len(kvParts) != 2 {
				err := fmt.Errorf("unexpected format: expect 2 kvParts, got '%v'", len(kvParts))
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			// splitting "table.ChannelTitle[0].Name"
			keyParts := strings.Split(kvParts[0], ".")
			switch len(keyParts) {
			case 3: // pass
			default:
				err := fmt.Errorf("unexpected format: expect 3 keyParts, got '%v'", len(keyParts))
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			// splitting "ChannelTitle[0]"
			chTitleParts := strings.Split(keyParts[1], "[")
			if len(chTitleParts) != 2 {
				err := fmt.Errorf("unexpected format: expect 2 chTitleParts, got '%v'", len(chTitleParts))
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			// trimming "0]"
			chanNumber := strings.TrimRight(chTitleParts[1], "]")
			if chanNumber == "" {
				err := fmt.Errorf("unexpected format: expect non empty chanNumber, got '%s'", chanNumber)
				ev.MarkWithGoError(err)
				ev.DumpThis(string(body[:]))
				return nil, ev
			}

			parsedList = append(parsedList, kvParts[1])
		}
		
		successResponse = parsedList

	default:
		failedResponse := &TXT_ResponseStatus{
			RequestURL: reqURL,
			StatusCode: resp.StatusCode,
			StatusMsg : string(body),
		}

		ev.MarkWithAPIError(failedResponse)
		ev.DumpThis(failedResponse.FullError())
		return nil, ev
	}

	return successResponse, ev
}
