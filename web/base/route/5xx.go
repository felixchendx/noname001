package route

import (
	"html/template"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	baseEmbedding "noname001/web/base/embedding"
)

func (rh *BaseRouteHandler) Route500(ctx *fasthttp.RequestCtx, err error) {
	rh.logger.Debugf("%s: 500 %7s %s, err: %s", rh.logPrefix, string(ctx.Method()), string(ctx.Path()), err.Error())

	// TODO: json support
	// TODO: more contextual message, at least message that's helpful for reporting
	switch {
	case rh.authProvider.IsLoggedIn(ctx):
		rh.render500Authenticated(ctx)
	default: 
		rh.render500(ctx)
	}
}

func (rh *BaseRouteHandler) render500Authenticated(ctx *fasthttp.RequestCtx) {
	pageData := &baseEmbedding.PageData_default{}
	pageData.Title = "500 - Internal Server Error"
	pageData.ContentData = map[string]any{}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--500.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.logger.Errorf("%s: render500Authenticated err %s", rh.logPrefix, string(ctx.Path()), renderErr.Error())
		renderOut = template.HTML(`
500 - internal server error.
Something went wrong. Our programmer will be working on it.
`)
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetBody([]byte(renderOut))
}

func (rh *BaseRouteHandler) render500(ctx *fasthttp.RequestCtx) {
	rendered, renderErr := rh.templating.RenderCommon("p-500.html.tmpl", map[string]any{})
	if renderErr != nil {
		rh.logger.Errorf("%s: render500 err %s", rh.logPrefix, string(ctx.Path()), renderErr.Error())
		rendered = template.HTML(`
500 - internal server error.
Something went wrong. Our programmer will be working on it.
`)
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetBody([]byte(rendered))
}