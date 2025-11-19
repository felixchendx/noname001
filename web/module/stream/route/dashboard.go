package route

import (
	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
)

func (rh *ModuleRouteHandler) renderDashboard(ctx *fasthttp.RequestCtx) {
	flash := rh.baseBundle.Flash.GetFlashBundle(ctx)

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Stream - Dashboard"
	pageData.Messages = flash.Prev.Messages
	pageData.ContentData = map[string]any{
		"_title": "Dashboard",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--dashboard.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
