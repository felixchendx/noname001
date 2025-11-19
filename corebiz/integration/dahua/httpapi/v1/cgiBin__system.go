package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"noname001/corebiz/integration/base/apicall"
)

func (api *APIClient) GetMachineName() (machineName string, apiEvRet *apicall.APICallEvent) {
	var uriApiDeviceName = "/cgi-bin/magicBox.cgi?action=getMachineName"
	ev := apicall.NewEvent("GetMachineName")
	defer ev.MarkAsEnded()

	reqCtx, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL+uriApiDeviceName,
		nil,
	)
	if err != nil {
		ev.MarkWithGoError(err)
		return "", ev
	}

	body, ev := api.sendRequestHttp(req, uriApiDeviceName, ev)
	if (ev.IsConsideredError()) || (body == "") {
		return "", ev
	}

	machineName, err = api.parseSingleValueString(body, "name")
	if err != nil {
		ev.MarkWithGoError(err)
		ev.DumpThis(string(body[:]))
		return "", ev
	}

	return machineName, ev
}

func (api *APIClient) GetDeviceType() (deviceType string, apiEvRet *apicall.APICallEvent) {
	var uriApiDeviceType = "/cgi-bin/magicBox.cgi?action=getDeviceType"
	ev := apicall.NewEvent("GetDeviceType")
	defer ev.MarkAsEnded()

	reqCtx, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL+uriApiDeviceType,
		nil,
	)
	if err != nil {
		ev.MarkWithGoError(err)
		return "", ev
	}

	body, ev := api.sendRequestHttp(req, uriApiDeviceType, ev)
	if (ev.IsConsideredError()) || (body == "") {
		return "", ev
	}

	deviceType, err = api.parseSingleValueString(body, "type")
	if err != nil {
		ev.MarkWithGoError(err)
		ev.DumpThis(string(body[:]))
		return "", ev
	}

	return deviceType, ev
}

func (api *APIClient) GetVendor() (manufacture string, apiEvRet *apicall.APICallEvent) {
	var uriApiManufacture = "/cgi-bin/magicBox.cgi?action=getVendor"
	ev := apicall.NewEvent("GetVendor")
	defer ev.MarkAsEnded()

	reqCtx, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL+uriApiManufacture,
		nil,
	)
	if err != nil {
		ev.MarkWithGoError(err)
		return "", ev
	}

	body, ev := api.sendRequestHttp(req, uriApiManufacture, ev)
	if (ev.IsConsideredError()) || (body == "") {
		return "", ev
	}

	manufacture, err = api.parseSingleValueString(body, "vendor")
	if err != nil {
		ev.MarkWithGoError(err)
		ev.DumpThis(string(body[:]))
		return "", ev
	}

	return manufacture, ev
}

func (api *APIClient) GetHardwareVersion() (HardwareVersion string, apiEvRet *apicall.APICallEvent) {
	var uriApiHardwareVersion = "/cgi-bin/magicBox.cgi?action=getHardwareVersion"
	ev := apicall.NewEvent("GetHardwareVersion")
	defer ev.MarkAsEnded()

	reqCtx, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL+uriApiHardwareVersion,
		nil,
	)
	if err != nil {
		errReturn := fmt.Errorf("%s-%s", err.Error())
		ev.MarkWithGoError(errReturn)
		return "", ev
	}

	body, ev := api.sendRequestHttp(req, uriApiHardwareVersion, ev)
	if (ev.IsConsideredError()) || (body == "") {
		return "", ev
	}

	HardwareVersion, err = api.parseSingleValueString(body, "version")
	if err != nil {
		ev.MarkWithGoError(err)
		ev.DumpThis(string(body[:]))
		return "", ev
	}

	return HardwareVersion, ev
}

func (api *APIClient) GetSoftwareVersion() (*TXT_SoftwareVersion, *apicall.APICallEvent) {
	var uriApiSoftwareVersion = "/cgi-bin/magicBox.cgi?action=getSoftwareVersion"
	ev := apicall.NewEvent("GetSoftwareVersion")
	defer ev.MarkAsEnded()

	reqCtx, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL+uriApiSoftwareVersion,
		nil,
	)
	if err != nil {
		ev.MarkWithGoError(err)
		return nil, ev
	}

	body, ev := api.sendRequestHttp(req, uriApiSoftwareVersion, ev)
	if (ev.IsConsideredError()) || (body == "") {
		return nil, ev
	}

	var softwareVersion *TXT_SoftwareVersion = &TXT_SoftwareVersion{}
	softwareVersion.Version, softwareVersion.ReleaseDate, err = api.parseSoftwareVersion(body)
	if err != nil {
		ev.MarkWithGoError(err)
		ev.DumpThis(string(body[:]))
		return nil, ev
	}

	return softwareVersion, ev
}

func (api *APIClient) GetSystemInfo() (*TXT_SystemInfo, *apicall.APICallEvent) {
	var uriApiSystemInfo = "/cgi-bin/magicBox.cgi?action=getSystemInfo"
	ev := apicall.NewEvent("GetSystemInfo")
	defer ev.MarkAsEnded()

	reqCtx, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL+uriApiSystemInfo,
		nil,
	)
	if err != nil {
		ev.MarkWithGoError(err)
		return nil, ev
	}

	body, ev := api.sendRequestHttp(req, uriApiSystemInfo, ev)
	if (ev.IsConsideredError()) || (body == "") {
		return nil, ev
	}

	var systemInfo *TXT_SystemInfo = &TXT_SystemInfo{}
	systemInfo.Processor, systemInfo.SerialNumber, systemInfo.UpdateSerial, err = api.parseSystemInfo(body)
	if err != nil {
		ev.MarkWithGoError(err)
		ev.DumpThis(string(body[:]))
		return nil, ev
	}

	return systemInfo, ev
}


func (api *APIClient) GetVideoEncodeConfig(channelID string) (*TXT_VideoEncodeConfig, *apicall.APICallEvent) {
	ev := apicall.NewEvent("GetVideoEncodeConfig")
	defer ev.MarkAsEnded()

	reqCtx, reqCancel := context.WithTimeout(api.context, api.httpTimeout)
	defer reqCancel()

	channelIDActualInt, err := strconv.Atoi(channelID)
	if err != nil {
		ev.MarkWithGoError(err)
		return nil, ev
	} else if (channelIDActualInt - 1) < 0 {
		ev.MarkWithGoError(fmt.Errorf("channel id does not valid"))
		return nil, ev
	}
	var ChannelIDActualStr = fmt.Sprintf("%d", channelIDActualInt-1)
	var uriVideoEncodeConfig = fmt.Sprintf("/cgi-bin/configManager.cgi?action=getConfig&name=Encode[%s]", ChannelIDActualStr) 

	req, err := http.NewRequestWithContext(
		reqCtx,
		"GET", api.baseURL+uriVideoEncodeConfig,
		nil,
	)
	if err != nil {
		ev.MarkWithGoError(err)
		return nil, ev
	}

	body, ev := api.sendRequestHttp(req, uriVideoEncodeConfig, ev)
	if (ev.IsConsideredError()) || (body == "") {
		return nil, ev
	}

	var videoConfig *TXT_VideoEncodeConfig = &TXT_VideoEncodeConfig{}
	videoConfig, err = api.parseVideoEncodeConfig(ChannelIDActualStr, body)
	if err != nil {
		ev.MarkWithGoError(err)
		ev.DumpThis(string(body[:]))
		return nil, ev
	}

	return videoConfig, ev
}

func (api *APIClient) sendRequestHttp(req *http.Request, uri string, ev *apicall.APICallEvent) (string, *apicall.APICallEvent) {
	resp, err := api.httpClient.Do(req)
	if err != nil {
		ev.MarkWithGoError(err)
		return "", ev
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ev.MarkWithGoError(err)
		return "", ev
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		default:
			failedResponse := &TXT_ResponseStatus{
				RequestURL: uri,
				StatusCode: resp.StatusCode,
				StatusMsg:  string(body),
			}
			ev.MarkWithAPIError(failedResponse)
			ev.DumpThis(failedResponse.FullError())
		}
		return "", ev
	}

	return string(body), ev
}

