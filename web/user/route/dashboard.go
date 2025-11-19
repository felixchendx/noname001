package route

import (
	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
)

func (rh *UserRouteHandler) renderDashboard(ctx *fasthttp.RequestCtx) {
	flash := rh.baseBundle.Flash.GetFlashBundle(ctx)

	webUser := rh.baseBundle.Auth.CurrentWebUser(ctx)

	pageData := rh.templating.NewPageData()
	pageData.Title = "User - Dashboard"
	pageData.Messages = flash.Prev.Messages
	pageData.ExtraCssLinks = []string{}
	pageData.ContentData = map[string]any{
		"_title": "Dashboard",

		"_username": webUser.Name,
		"_userrole": webUser.Role,
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
