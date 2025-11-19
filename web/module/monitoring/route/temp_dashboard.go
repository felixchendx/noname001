package route

import (
	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
)

// this whole temp dashboard is temporary
// both as minimum simple dumb monitoring
// and more importantly, TO TEST USING CACHE STUFFS

// anw, full fledge monitoring is lotta works, gotta work on higher priority stuffs

func (rh *ModuleRouteHandler) renderTempDashboard(ctx *fasthttp.RequestCtx) {
	flash := rh.baseBundle.Flash.GetFlashBundle(ctx)

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Monitoring - Temp Dashboard"
	pageData.Messages = flash.Prev.Messages
	pageData.ExtraCssLinks = []string {
		"/monitoring/assets/temp-dashboard.css",
	}
	pageData.ExtraJsLinks = []string{
		"/assets/internal/ws.js",
		"/monitoring/assets/temp-dashboard.js",
	}
	pageData.ContentData = map[string]any{
		"_title": "Temp Dashboard",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--temp-dashboard.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
