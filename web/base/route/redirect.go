package route

import (
	"github.com/valyala/fasthttp"
)

func (rh *BaseRouteHandler) RedirectTo(ctx *fasthttp.RequestCtx, uri []byte) {
	rh.logger.Debugf("%s: 302 %7s %s -> %s", rh.logPrefix, string(ctx.Method()), string(ctx.Path()), string(uri))
	ctx.RedirectBytes(uri, fasthttp.StatusFound)
}

func (rh *BaseRouteHandler) RedirectToIndex(ctx *fasthttp.RequestCtx) {
	ctx.RedirectBytes([]byte("/"), fasthttp.StatusFound)
}

func (rh *BaseRouteHandler) RedirectToLogin(ctx *fasthttp.RequestCtx) {
	ctx.RedirectBytes([]byte("/login"), fasthttp.StatusFound)
}
