package route

import (
	"bytes"

	"github.com/valyala/fasthttp"

	"noname001/logging"
	webBase "noname001/web/base"
	"noname001/web/module/device/embedding"
)

type ModuleRouteHandlerParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle

	Templating *embedding.ModuleTemplating
}

type ModuleRouteHandler struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	templating   *embedding.ModuleTemplating

	pathPrefixes [][]byte

	assetHandler fasthttp.RequestHandler
}

func NewModuleRouteHandler(params *ModuleRouteHandlerParams) *ModuleRouteHandler {
	rh := &ModuleRouteHandler{}
	rh.logger = params.Logger
	rh.logPrefix = params.LogPrefix + ".rh"

	rh.baseBundle = params.BaseBundle

	rh.templating = params.Templating

	rh.pathPrefixes = [][]byte{
		[]byte("/device"),
	}

	rh.assetHandler = embedding.FakeSingletonAssetHandler()

	return rh
}

func (rh *ModuleRouteHandler) PathPrefixes() [][]byte {
	return rh.pathPrefixes
}

func (rh *ModuleRouteHandler) HasPrefix(path []byte) bool {
	for _, pathPrefix := range rh.pathPrefixes {
		if bytes.HasPrefix(path, pathPrefix) {
			return true
		}
	}

	return false
}

func (rh *ModuleRouteHandler) RootHandler(ctx *fasthttp.RequestCtx) {
	switch {
	case rh.baseBundle.Navi.DoesNotHave("device"): rh.baseBundle.RouteHandler.Route404_ModuleNotActive(ctx, "device"); return
	case rh.baseBundle.Auth.IsLoggedOut(ctx)     : rh.baseBundle.RouteHandler.RedirectToLogin(ctx); return
	}

	path := bytes.TrimRight(ctx.Path(), "/")
	isGET := ctx.IsGet()
	isPOST := ctx.IsPost()
	___redirectTo := rh.baseBundle.RouteHandler.RedirectTo

	switch {
	case isGET && bytes.Equal(path, []byte("/device"))               : ___redirectTo(ctx, []byte("/device/device/listing"))
	case isGET && bytes.Equal(path, []byte("/device/dashboard"))     : rh.renderDashboard(ctx)

	case isGET && bytes.Equal(path, []byte("/device/device"))        : ___redirectTo(ctx, []byte("/device/device/listing"))
	case isGET && bytes.Equal(path, []byte("/device/device/listing")): rh.renderDeviceListing(ctx)

	case isGET  && bytes.Equal(path, []byte("/device/device/detail"))          : rh.renderDeviceDetail(ctx)
	case isGET  && bytes.Equal(path, []byte("/device/device/detail/ws"))       : rh.deviceDetail__ws(ctx)
	case isPOST && bytes.Equal(path, []byte("/device/device/detail/do/add"))   : rh.doAddDevice(ctx)
	case isPOST && bytes.Equal(path, []byte("/device/device/detail/do/edit"))  : rh.doEditDevice(ctx)
	case isPOST && bytes.Equal(path, []byte("/device/device/detail/do/delete")): rh.doDeleteDevice(ctx)

	case bytes.HasPrefix(path, []byte("/device/assets")): rh.assetHandler(ctx)

	default:
		rh.baseBundle.RouteHandler.Route404(ctx)
	}

}
