package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
	webUtil     "noname001/web/base/util"

	"noname001/app/sys"
)

func (rh *AdminRouteHandler) renderUserListing(ctx *fasthttp.RequestCtx) {
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
		case "username", "role_simple" : dataFilter[sk] = sv
		}
	})

	sc := &sys.SysUser__SearchCriteria{
		UsernameLike: dataFilter["username"],
		RoleSimple: []string{dataFilter["role_simple"]},

		Pagination: &sys.SearchPagination{
			CurrPage: currPage,
			PerPage: perPage,
		},
	}
	searchResult, searchMessages := sys.Bundle.Service.SysUser__Search(sc)
	showingMessages.Append(searchMessages)

	if searchResult != nil {
		for _, item := range searchResult.Data {
			dataRows = append(dataRows, map[string]string{
				"username": item.Username,
				"role_simple": item.RoleSimple,

				"link_detail": fmt.Sprintf("%s?id=%s", "/admin/user/detail", item.ID),
			})
		}

		dataPg = rh.baseBundle.Util.GenerateWebPagination(&webUtil.WebPaginationParams{
			SearchPagination: searchResult.Pagination,
			LinkURI: "/admin/user/listing",
			QueryParams: dataFilter,
			Anchor: "datatable",
		})
	}

	pageData := rh.templating.NewPageData()
	pageData.Title = "Admin - User listing"
	pageData.Messages = showingMessages
	pageData.ContentData = map[string]any{
		"_title": "User listing",
		"_link": map[string]string{
			"user_listing": "/admin/user/listing",
			"new_user": "/admin/user/detail",
		},
		"_datatable": map[string]any{
			"user_listing": map[string]any{
				"filter": dataFilter,
				"rows": dataRows,
				"pg": dataPg,
			},
		},
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--user-listing.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
