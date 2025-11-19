package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
	webUtil     "noname001/web/base/util"

	wallService "noname001/app/module/feature/wall/service"
)

func (rh *ModuleRouteHandler) renderWallListing(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages = flash.Prev.Messages

		dataFilter map[string]string      = make(map[string]string)
		dataRows   []map[string]string    = make([]map[string]string, 0)
		dataPg     *webUtil.WebPagination
	)

	currPage, perPage := rh.baseBundle.Util.ExtractPaginationFromQueryParams(ctx)

	ctx.QueryArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "code", "name", "state" : dataFilter[sk] = sv
		default: // ignored
		}
	})

	sc := &wallService.Web__WallListingSearchCriteria{
		WallCodeLike: dataFilter["code"],
		WallNameLike: dataFilter["name"],
		WallState: dataFilter["state"],
		// WallLayout

		Pagination: &wallService.SearchPagination{
			CurrPage: currPage,
			PerPage: perPage,
		},
	}
	searchResult, searchMessages := wallService.Instance().Web__WallListing(sc)
	showingMessages.Append(searchMessages)

	if searchResult != nil {
		for _, item := range searchResult.Data {
			item["link_detail"] = fmt.Sprintf("%s?id=%s", "/wall/wall/detail", item["wall_id"])
			item["link_view"]   = fmt.Sprintf("%s?id=%s", "/wall/wall/view", item["wall_id"])
			dataRows = append(dataRows, item)
		}

		dataPg = rh.baseBundle.Util.GenerateWebPagination(&webUtil.WebPaginationParams{
			SearchPagination: searchResult.Pagination,
			LinkURI: "/wall/wall/listing",
			QueryParams: dataFilter,
			Anchor: "datatable",
		})
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Wall - Wall listing"
	pageData.Messages = showingMessages
	pageData.ContentData = map[string]any{
		"_title": "Wall listing",
		"_link": map[string]string{
			"wall_listing": "/wall/wall/listing",
			"new_wall": "/wall/wall/detail",
		},
		"_datatable": map[string]any{
			"wall_listing": map[string]any{
				"filter": dataFilter,
				"rows": dataRows,
				"pg": dataPg,
			},
		},
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--wall-listing.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
