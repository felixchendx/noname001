package route

import (
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"

	webUtil "noname001/web/base/util"
	webConstant "noname001/web/constant"

	streamService "noname001/app/module/feature/stream/service"
)

func (rh *ModuleRouteHandler) renderStreamProfileListing(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages = flash.Prev.Messages

		dataFilter map[string]string      = make(map[string]string)
		dataRows   []map[string]string    = make([]map[string]string, 0)
		dataPg     *webUtil.WebPagination
	)

	currPage, perPage := rh.baseBundle.Util.ExtractPaginationFromQueryParams(ctx)

	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "code", "name", "state": dataFilter[sKey] = sValue
		}
	})

	sc := &streamService.StreamProfile__SearchCriteria{
		CodeLike: (dataFilter["code"]),
		NameLike: dataFilter["name"],
		State:    []string{dataFilter["state"]},

		Pagination: &streamService.SearchPagination{
			CurrPage: currPage,
			PerPage:  perPage,
		},
	}
	streamProfileSR, searchMessages := streamService.Instance().SearchStreamProfile(sc)
	showingMessages.Append(searchMessages)

	if streamProfileSR != nil {
		for _, item := range streamProfileSR.Data {
			var itemStreamProfile map[string]string = make(map[string]string)
			itemStreamProfile["id"]    = item.ID
			itemStreamProfile["code"]  = item.Code
			itemStreamProfile["name"]  = item.Name
			itemStreamProfile["state"] = item.State
			itemStreamProfile["note"]  = item.Note

			itemStreamProfile["target_video_codec"]       = item.TargetVideoCodec
			itemStreamProfile["target_video_compression"] = strconv.Itoa(item.TargetVideoCompression)
			itemStreamProfile["target_video_bitrate"]     = strconv.Itoa(item.TargetVideoBitrate)

			itemStreamProfile["target_audio_codec"]       = item.TargetAudioCodec
			itemStreamProfile["target_audio_compression"] = strconv.Itoa(item.TargetAudioCompression)
			itemStreamProfile["target_audio_bitrate"]     = strconv.Itoa(item.TargetAudioBitrate)

			itemStreamProfile["link_detail"] = fmt.Sprintf("%s?id=%s", "/stream/stream-profile/detail", item.ID)

			dataRows = append(dataRows, itemStreamProfile)
		}

		dataPg = rh.baseBundle.Util.GenerateWebPagination(&webUtil.WebPaginationParams{
			SearchPagination: streamProfileSR.Pagination,
			LinkURI:          "/stream/stream-profile/listing",
			QueryParams:      dataFilter,
			Anchor:           "datatable",
		})
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Stream - Stream Profile Listing"
	pageData.Messages = showingMessages
	pageData.ContentData = map[string]any{
		"_title": "Stream Profile Listing",
		"_link": map[string]any{
			"sp_listing": "/stream/stream-profile/listing",
			"new_sp":     "/stream/stream-profile/detail",
		},
		"_datatable": map[string]any{
			"sp_listing": map[string]any{
				"filter": dataFilter,
				"rows":   dataRows,
				"pg":     dataPg,
			},
		},
	}

	renderOut, err := rh.templating.RenderContent_default("d-content--stream-profile-listing.html.tmpl", pageData, ctx)
	if err != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, err)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
