package route

import (
	"strings"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	wallService "noname001/app/module/feature/wall/service"
)

func (rh *ModuleRouteHandler) renderWallView(ctx *fasthttp.RequestCtx) {
	var (
		wallServiceInstance = wallService.Instance()

		host     []byte   = ctx.Request.Header.Peek("Host")
		hostPart []string = strings.Split(string(host), ":")
	)

	dataID := string(ctx.QueryArgs().Peek("id"))

	flash := rh.baseBundle.Flash.GetFlashBundle(ctx)

	wall, wallMessages := wallServiceInstance.Wall__Find(dataID, true)
	if wallMessages.HasError() {
		flash.Next.Messages.Append(wallMessages)
		ctx.Redirect("/wall/wall", fasthttp.StatusFound)
		return
	}

	for _, wallItem := range wall.Items {
		wallItem.TempRelayURL = wallServiceInstance.GetRelayedStreamViewURL(hostPart[0], wallItem.SourceNodeID, wallItem.StreamCode, "hls")
	}

	wallLayout, wallLayoutMessages := wallServiceInstance.WallLayout__Find(wall.WallLayoutID)
	if wallLayoutMessages.HasError() {
		flash.Next.Messages.Append(wallLayoutMessages)
		ctx.Redirect("/wall/wall", fasthttp.StatusFound)
		return
	}

	contentTemplateToUse := ""
	switch wallLayout.Code {
	case "DEFAULT_4"   : contentTemplateToUse = "d-content--wall-view-default-4-topbar.html.tmpl"
	case "DEFAULT_12"  : contentTemplateToUse = "d-content--wall-view-default-12-topbar.html.tmpl"
	case "DEFAULT_16"  : contentTemplateToUse = "d-content--wall-view-default-16-topbar.html.tmpl"
	case "DEFAULT_1B7S": contentTemplateToUse = "d-content--wall-view-default-1b7s-topbar.html.tmpl"
	default            : contentTemplateToUse = "d-content--wall-view-default-4-topbar.html.tmpl"
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Wall - Wall View"
	pageData.Messages = flash.Prev.Messages
	pageData.Theme = webConstant.TBLR_THEME_DARK
	pageData.RenderSidenav = false
	pageData.RenderTopnav = false
	pageData.RenderFooter = false
	pageData.ExtraCssLinks = []string {
		"/assets/fu-component/live-stream/default/index.css",
		"/wall/assets/wall-view.css",
	}
	pageData.ExtraJsLinks = []string{
		"/assets/hls/hls.js",
		"/assets/internal/ws.js",
		"/assets/fu-component/live-stream/default/index.js",
		"/wall/assets/wall-view.js",
	}
	pageData.ContentData = map[string]any{
		"wall": wall,
	}

	renderOut, renderErr := rh.templating.RenderContent_default(contentTemplateToUse, pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}
