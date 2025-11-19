package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webUtil "noname001/web/base/util"
	webConstant "noname001/web/constant"

	deviceService "noname001/app/module/common/device/service"
)

func (rh *ModuleRouteHandler) renderDeviceListing(ctx *fasthttp.RequestCtx) {
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
		case "code", "name", "state", "hostname", "username", "brand":
			dataFilter[sKey] = sValue
		}
	})

	searchCriteria := deviceService.Web__DeviceListingSearchCriteria{
		DeviceCodeLike:     dataFilter["code"],
		DeviceNameLike:     dataFilter["name"],
		DeviceState:        dataFilter["state"],
		DeviceHostnameLike: dataFilter["hostname"],
		DeviceUsernameLike: dataFilter["username"],
		DeviceBrand:    dataFilter["brand"],

		Pagination: &deviceService.SearchPagination{
			CurrPage: currPage,
			PerPage:  perPage,
		},
	}

	searchResult, searchMessages := deviceService.Instance().Web__DeviceListing(&searchCriteria)
	showingMessages.Append(searchMessages)

	if searchResult != nil {
		for _, item := range searchResult.Data {
			item["link_detail"] = fmt.Sprintf("%s?id=%s", "/device/device/detail", item["device_id"])
			dataRows = append(dataRows, item)
		}

		dataPg = rh.baseBundle.Util.GenerateWebPagination(&webUtil.WebPaginationParams{
			SearchPagination: searchResult.Pagination,
			LinkURI:          "/device/device/listing",
			QueryParams:      dataFilter,
			Anchor:           "datatable",
		})
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Device - Device listing"
	pageData.Messages = showingMessages
	pageData.ContentData = map[string]any{
		"_title": "Device Listing",
		"_link": map[string]string{
			"device_listing": "/device/device/listing",
			"new_device":     "/device/device/detail",
		},
		"_datatable": map[string]any{
			"device_listing": map[string]any{
				"filter": dataFilter,
				"rows":   dataRows,
				"pg":     dataPg,
			},
		},
	}
	renderOut, renderErr := rh.templating.RenderContent_default("d-content--device-listing.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
