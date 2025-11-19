package route

import (
	"github.com/valyala/fasthttp"
)

func (rh *BaseRouteHandler) RouteMaintenance(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	// TODO: render https://preview.tabler.io/error-maintenance.html
	ctx.SetBody([]byte("TODO: render maintenance"))
}
