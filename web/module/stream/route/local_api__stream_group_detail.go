package route

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/valyala/fasthttp"

	webUtil "noname001/web/base/util"
	webConstant "noname001/web/constant"

	streamService "noname001/app/module/feature/stream/service"
)

type localApi_streamItem struct {
	ID               string `json:"id"`
	StreamGroupID    string `json:"stream_group_id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	State            string `json:"state"`
	Note             string `json:"note"`

	SourceType       string `json:"source_type"`
	DeviceCode       string `json:"device_code"`
	DeviceChannelID  string `json:"device_channel_id"`
	DeviceStreamType string `json:"device_stream_type"`
	ExternalURL      string `json:"external_url"`
	Filepath         string `json:"filepath"`
	EmbeddedFilepath string `json:"embedded_filepath"`
}

func (ls *localApi_streamItem) toDE() *streamService.StreamItemDE {
	return &streamService.StreamItemDE{
		ID:               ls.ID,
		Code:             ls.Code,
		Name:             ls.Name,
		State:            ls.State,
		Note:             ls.Note,
		StreamGroupID:    ls.StreamGroupID,
		SourceType:       ls.SourceType,
		DeviceCode:       ls.DeviceCode,
		DeviceChannelID:  ls.DeviceChannelID,
		DeviceStreamType: ls.DeviceStreamType,
		ExternalURL:      ls.ExternalURL,
		Filepath:         ls.Filepath,
		EmbeddedFilepath: ls.EmbeddedFilepath,
	}
}

func (rh *ModuleRouteHandler) localAPI_addStreamItem(ctx *fasthttp.RequestCtx) {
	reqBody := &localApi_streamItem{}
	err := json.Unmarshal(ctx.PostBody(), reqBody)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	de := reqBody.toDE()

	var defResp *webUtil.LAPIDefaultResponse

	_, messages := streamService.Instance().AddStreamItem(de)
	if messages.HasError() {
		defResp = rh.baseBundle.Util.OnelinerErrorLAPIResponse(nil, messages.FirstErrorMessageString())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		defResp = rh.baseBundle.Util.OnelinerOkLAPIResponse(nil, messages.FirstNoticeMessageString())
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	respBytes, err := json.Marshal(defResp)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_JSON)
	ctx.SetBody(respBytes)
}


func (rh *ModuleRouteHandler) localAPI_editStreamItem(ctx *fasthttp.RequestCtx) {
	reqBody := &localApi_streamItem{}
	err := json.Unmarshal(ctx.PostBody(), reqBody)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	de := reqBody.toDE()

	var defResp *webUtil.LAPIDefaultResponse

	_, messages := streamService.Instance().EditStreamItem(reqBody.ID, de)
	if messages.HasError() {
		defResp = rh.baseBundle.Util.OnelinerErrorLAPIResponse(nil, messages.FirstErrorMessageString())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		defResp = rh.baseBundle.Util.OnelinerOkLAPIResponse(nil, messages.FirstNoticeMessageString())
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	respBytes, err := json.Marshal(defResp)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_JSON)
	ctx.SetBody(respBytes)
}

func (rh *ModuleRouteHandler) localAPI_deleteStreamItem(ctx *fasthttp.RequestCtx) {
	reqBody := &localApi_streamItem{}
	err := json.Unmarshal(ctx.PostBody(), reqBody)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	de := streamService.Instance().EmptyStreamItem()
	de.ID = reqBody.ID

	var defResp *webUtil.LAPIDefaultResponse

	messages := streamService.Instance().DeleteStreamItem(de.ID)
	if messages.HasError() {
		defResp = rh.baseBundle.Util.OnelinerErrorLAPIResponse(nil, messages.FirstErrorMessageString())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		defResp = rh.baseBundle.Util.OnelinerOkLAPIResponse(nil, messages.FirstNoticeMessageString())
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	respBytes, err := json.Marshal(defResp)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_JSON)
	ctx.SetBody(respBytes)
}

type localAPI__deviceChannelPreview struct {
	DeviceCode       string `json:"device_code"`	
	DeviceChannelID  string `json:"device_channel_id"`
	DeviceStreamType string `json:"device_stream_type"`

	PreviewURL       string `json:"preview_url"`
}
func (rh *ModuleRouteHandler) localAPI__requestDeviceChannelPreview(ctx *fasthttp.RequestCtx) {
	host := ctx.Request.Header.Peek("Host")
	hostPart := strings.Split(string(host), ":")

	reqBody := &localAPI__deviceChannelPreview{}
	err := json.Unmarshal(ctx.PostBody(), reqBody)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	var defResp *webUtil.LAPIDefaultResponse
	previewURL, messages := streamService.Instance().DeviceChannelPreview(hostPart[0], reqBody.DeviceCode, reqBody.DeviceChannelID, reqBody.DeviceStreamType, "hls")
	if messages.HasError() {
		defResp = rh.baseBundle.Util.OnelinerErrorLAPIResponse(nil, messages.LastErrorMessageString())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		reqBody.PreviewURL = previewURL
		time.Sleep(3 * time.Second) // temp artificial delay for ffmpeg -> mediamtx

		defResp = rh.baseBundle.Util.OnelinerOkLAPIResponse(reqBody, "")
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	respBytes, err := json.Marshal(defResp)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_JSON)
	ctx.SetBody(respBytes)
}
