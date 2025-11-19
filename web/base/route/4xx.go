package route

import (
	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	baseEmbedding "noname001/web/base/embedding"
)

func (rh *BaseRouteHandler) Route403(ctx *fasthttp.RequestCtx) {
	rh.logger.Debugf("%s: 403 %7s %s", rh.logPrefix, string(ctx.Method()), string(ctx.Path()))

	// TODO: json support
	switch {
	case rh.authProvider.IsLoggedIn(ctx):
		rh.render403Authenticated(ctx)
	default: 
		rh.render403(ctx)
	}
}

func (rh *BaseRouteHandler) Route404(ctx *fasthttp.RequestCtx) {
	rh.logger.Debugf("%s: 404 %7s %s", rh.logPrefix, string(ctx.Method()), string(ctx.Path()))

	// TODO: check request header content type to support json response
	switch {
	case rh.authProvider.IsLoggedIn(ctx):
		rh.render404Authenticated(ctx)
	default: 
		rh.render404(ctx)
	}
}

func (rh *BaseRouteHandler) Route404_ModuleNotActive(ctx *fasthttp.RequestCtx, moduleCode string) {
	rh.logger.Debugf("%s: 404 %7s %s", rh.logPrefix, string(ctx.Method()), string(ctx.Path()))

	rh.render404ModuleNotActive(ctx, moduleCode)
}


func (rh *BaseRouteHandler) render403Authenticated(ctx *fasthttp.RequestCtx) {
	pageData := &baseEmbedding.PageData_default{}
	pageData.Title = "403 - Forbidden"
	pageData.ContentData = map[string]any{}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--403-authenticated.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusForbidden)
	ctx.SetBody([]byte(renderOut))
}

func (rh *BaseRouteHandler) render403(ctx *fasthttp.RequestCtx) {
	rendered, renderErr := rh.templating.RenderCommon("p-403.html.tmpl", map[string]any{})
	if renderErr != nil {
		rh.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusForbidden)
	ctx.SetBody([]byte(rendered))
}

func (rh *BaseRouteHandler) render404Authenticated(ctx *fasthttp.RequestCtx) {
	pageData := &baseEmbedding.PageData_default{}
	pageData.Title = "404 - Page not found"
	pageData.ContentData = map[string]any{}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--404-authenticated.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBody([]byte(renderOut))
}

func (rh *BaseRouteHandler) render404(ctx *fasthttp.RequestCtx) {
	rendered, renderErr := rh.templating.RenderCommon("p-404.html.tmpl", map[string]any{})
	if renderErr != nil {
		rh.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBody([]byte(rendered))
}

func (rh *BaseRouteHandler) render404ModuleNotActive(ctx *fasthttp.RequestCtx, moduleCode string) {
	pageData := &baseEmbedding.PageData_default{}
	pageData.Title = "404 - Page not found"
	pageData.ContentData = map[string]any{
		"_module_code": moduleCode,
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--404-module-not-active.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBody([]byte(renderOut))
}
