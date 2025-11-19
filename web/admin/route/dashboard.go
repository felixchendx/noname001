package route

import (
	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
)

// func (rh *AdminRouteHandler) routeDashboard(ctx *fasthttp.RequestCtx) {
// 	switch string(ctx.Method()) {
// 	case fasthttp.MethodGet:
// 		rh.renderDashboard(ctx)
// 	default:
// 		rh.baseBundle.RouteHandler.Route404(ctx)
// 	}

// 	// post exec shenanigans
// }

func (rh *AdminRouteHandler) renderDashboard(ctx *fasthttp.RequestCtx) {
	flash := rh.baseBundle.Flash.GetFlashBundle(ctx)

	pageData := rh.templating.NewPageData()
	pageData.Title = "Admin - Dashboard"
	pageData.Messages = flash.Prev.Messages
	pageData.ExtraCssLinks = []string{}
	pageData.ContentData = map[string]any{
		"_title": "Dashboard",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content-dashboard.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
