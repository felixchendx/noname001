package route

import (
	"github.com/valyala/fasthttp"
)

func (rh *BaseRouteHandler) RouteLogout(ctx *fasthttp.RequestCtx) {
	if rh.authProvider.IsLoggedOut(ctx) {
		rh.RedirectToLogin(ctx)
		return
	}

	rh.doLogout(ctx)
}

func (rh *BaseRouteHandler) doLogout(ctx *fasthttp.RequestCtx) {
	rh.authProvider.WebLogout(ctx)
	rh.RedirectToLogin(ctx)
}
