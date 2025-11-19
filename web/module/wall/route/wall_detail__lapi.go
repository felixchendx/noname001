package route

import (
	"encoding/json"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
	webUtil     "noname001/web/base/util"

	wallService "noname001/app/module/feature/wall/service"
)

type localAPI__updateWallItem struct {
	WallItemID   string `json:"wall_item_id"`
	SourceNodeID string `json:"source_node_id"`
	StreamCode   string `json:"stream_code"`
}
func (rh *ModuleRouteHandler) localAPI__updateWallItem(ctx *fasthttp.RequestCtx) {
	reqBody := &localAPI__updateWallItem{}
	err := json.Unmarshal(ctx.PostBody(), reqBody)
	if err != nil {
		rh.baseBundle.Util.Oneliner500LAPIResponse(ctx, err)
		return
	}

	de := wallService.Instance().WallItem__Empty()
	de.SourceNodeID = reqBody.SourceNodeID
	de.StreamCode = reqBody.StreamCode


	var defResp *webUtil.LAPIDefaultResponse

	_, messages := wallService.Instance().WallItem__Edit(reqBody.WallItemID, de)
	if messages.HasError() {
		defResp = rh.baseBundle.Util.OnelinerErrorLAPIResponse(messages.FirstErrorMessageString(), "")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		defResp = rh.baseBundle.Util.OnelinerOkLAPIResponse(nil, "")
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
