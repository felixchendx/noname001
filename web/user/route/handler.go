package route

import (
	"bytes"

	"github.com/valyala/fasthttp"

	"noname001/logging"

	webBase "noname001/web/base"
	"noname001/web/user/embedding"
)

type UserRouteHandlerParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle

	Templating *embedding.UserTemplating
}
type UserRouteHandler struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	templating   *embedding.UserTemplating

	pathPrefixes [][]byte

	assetHandler fasthttp.RequestHandler
}

func NewUserRouteHandler(params *UserRouteHandlerParams) (*UserRouteHandler) {
	rh := &UserRouteHandler{}
	rh.logger = params.Logger
	rh.logPrefix = params.LogPrefix + ".rh"

	rh.baseBundle = params.BaseBundle

	rh.templating = params.Templating

	rh.pathPrefixes = [][]byte{
		[]byte("/user"),
	}

	rh.assetHandler = embedding.FakeSingletonAssetHandler()

	return rh
}

func (rh *UserRouteHandler) PathPrefixes() ([][]byte) {
	return rh.pathPrefixes
}

func (rh *UserRouteHandler) HasPrefix(path []byte) (bool) {
	for _, pathPrefix := range rh.pathPrefixes {
		if bytes.HasPrefix(path, pathPrefix) {
			return true
		}
	}

	return false
}

func (rh *UserRouteHandler) RootHandler(ctx *fasthttp.RequestCtx) {
	switch {
	case rh.baseBundle.Auth.IsLoggedOut(ctx): rh.baseBundle.RouteHandler.RedirectToLogin(ctx); return
	}

	path   := bytes.TrimRight(ctx.Path(), "/")
	isGET  := ctx.IsGet()
	// isPOST := ctx.IsPost()

	___redirectTo := rh.baseBundle.RouteHandler.RedirectTo

	switch {
	case isGET && bytes.Equal(path, []byte("/user"))          : ___redirectTo(ctx, []byte("/user/dashboard"))
	case isGET && bytes.Equal(path, []byte("/user/dashboard")): rh.renderDashboard(ctx)

	case bytes.HasPrefix(path, []byte("/user/assets")): rh.assetHandler(ctx)

	default:
		rh.baseBundle.RouteHandler.Route404(ctx)
	}
}
