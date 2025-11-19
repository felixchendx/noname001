package route

import (
	"bytes"

	"github.com/valyala/fasthttp"

	"noname001/logging"

	webBase "noname001/web/base"

	"noname001/web/module/wall/embedding"
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

func NewModuleRouteHandler(params *ModuleRouteHandlerParams) (*ModuleRouteHandler) {
	rh := &ModuleRouteHandler{}
	rh.logger = params.Logger
	rh.logPrefix = params.LogPrefix + ".rh"

	rh.baseBundle = params.BaseBundle

	rh.templating = params.Templating

	rh.pathPrefixes = [][]byte{
		[]byte("/wall"),
	}

	rh.assetHandler = embedding.FakeSingletonAssetHandler()

	return rh
}

func (rh *ModuleRouteHandler) PathPrefixes() ([][]byte) {
	return rh.pathPrefixes
}

func (rh *ModuleRouteHandler) HasPrefix(path []byte) (bool) {
	for _, pathPrefix := range rh.pathPrefixes {
		if bytes.HasPrefix(path, pathPrefix) {
			return true
		}
	}
	return false
}

func (rh *ModuleRouteHandler) RootHandler(ctx *fasthttp.RequestCtx) {
	switch {
	case rh.baseBundle.Navi.DoesNotHave("wall"): rh.baseBundle.RouteHandler.Route404_ModuleNotActive(ctx, "wall"); return
	case rh.baseBundle.Auth.IsLoggedOut(ctx)   : rh.baseBundle.RouteHandler.RedirectToLogin(ctx); return
	}

	path   := bytes.TrimRight(ctx.Path(), "/")
	isGET  := ctx.IsGet()
	isPOST := ctx.IsPost()

	___redirectTo := rh.baseBundle.RouteHandler.RedirectTo

	switch {
	case isGET && bytes.Equal(path, []byte("/wall"))             : ___redirectTo(ctx, []byte("/wall/wall/listing"))
	case isGET && bytes.Equal(path, []byte("/wall/dashboard"))   : rh.renderDashboard(ctx)

	case isGET && bytes.Equal(path, []byte("/wall/wall"))        : ___redirectTo(ctx, []byte("/wall/wall/listing"))
	case isGET && bytes.Equal(path, []byte("/wall/wall/listing")): rh.renderWallListing(ctx)
	
	case isGET  && bytes.Equal(path, []byte("/wall/wall/detail"))          : rh.renderWallDetail(ctx)
	case isPOST && bytes.Equal(path, []byte("/wall/wall/detail/do/add"))   : rh.doAddWall(ctx)
	case isPOST && bytes.Equal(path, []byte("/wall/wall/detail/do/edit"))  : rh.doEditWall(ctx)
	case isPOST && bytes.Equal(path, []byte("/wall/wall/detail/do/delete")): rh.doDeleteWall(ctx)
	case isGET  && bytes.Equal(path, []byte("/wall/wall/detail/ws"))       : rh.wallDetail__ws(ctx)
	
	case isGET && bytes.Equal(path, []byte("/wall/wall/view"))   : rh.renderWallView(ctx)
	case isGET && bytes.Equal(path, []byte("/wall/wall/view/ws")): rh.wallView__ws(ctx)

	// TODO: reconfigure path, also restructure lapi stuffs
	case isPOST && bytes.Equal(path, []byte("/wall/local-api/update-wall-item")): rh.localAPI__updateWallItem(ctx)

	case bytes.HasPrefix(path, []byte("/wall/assets")): rh.assetHandler(ctx)

	default:
		rh.baseBundle.RouteHandler.Route404(ctx)
	}
}
