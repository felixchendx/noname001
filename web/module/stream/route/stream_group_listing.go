package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webUtil "noname001/web/base/util"
	webConstant "noname001/web/constant"

	streamService "noname001/app/module/feature/stream/service"
)

func (rh *ModuleRouteHandler) renderStreamGroupListing(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages = flash.Prev.Messages

		dataFilter map[string]string   = make(map[string]string)
		dataRows   []map[string]string = make([]map[string]string, 0)
		dataPg     *webUtil.WebPagination
	)

	currPage, perPage := rh.baseBundle.Util.ExtractPaginationFromQueryParams(ctx)

	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "code", "name", "state":
			dataFilter[sKey] = sValue
		}
	})

	sc := &streamService.Web__StreamGroupListingSearchCriteria{
		CodeLike: dataFilter["code"],
		NameLike: dataFilter["name"],
		State   : []string{dataFilter["state"]},

		Pagination: &streamService.SearchPagination{
			CurrPage: currPage,
			PerPage:  perPage,
		},
	}
	searchResult, searchMessages := streamService.Instance().Web__StreamGroupListing(sc)
	showingMessages.Append(searchMessages)

	if searchResult != nil {
		for _, item := range searchResult.Data {
			item["link_detail"] = fmt.Sprintf("%s?id=%s", "/stream/stream-group/detail-02", item["sg_id"])
			dataRows = append(dataRows, item)
		}

		dataPg = rh.baseBundle.Util.GenerateWebPagination(&webUtil.WebPaginationParams{
			SearchPagination: searchResult.Pagination,
			LinkURI         : "/stream/stream-group/listing",
			QueryParams     : dataFilter,
			Anchor          : "datatable",
		})
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Stream - Stream Group Listing"
	pageData.Messages = showingMessages
	pageData.ContentData = map[string]any{
		"_title": "Stream group listing",
		"_link": map[string]string{
			"sg_listing": "/stream/stream-group/listing",
			"new_sg":     "/stream/stream-group/detail-02",
		},
		"_datatable": map[string]any{
			"sg_listing": map[string]any{
				"filter": dataFilter,
				"rows":   dataRows,
				"pg":     dataPg,
			},
		},
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--stream-group-listing.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
