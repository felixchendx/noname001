package route

import (
	"bytes"

	"github.com/valyala/fasthttp"

	"noname001/logging"

	webBase "noname001/web/base"

	"noname001/web/admin/embedding"
)

type AdminRouteHandlerParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle

	Templating *embedding.AdminTemplating
}
type AdminRouteHandler struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	templating   *embedding.AdminTemplating

	pathPrefixes [][]byte

	assetHandler fasthttp.RequestHandler
}

func NewAdminRouteHandler(params *AdminRouteHandlerParams) (*AdminRouteHandler) {
	rh := &AdminRouteHandler{}
	rh.logger = params.Logger
	rh.logPrefix = params.LogPrefix + ".rh"

	rh.baseBundle = params.BaseBundle

	rh.templating = params.Templating

	rh.pathPrefixes = [][]byte{
		[]byte("/admin"),
	}

	rh.assetHandler = embedding.FakeSingletonAssetHandler()

	return rh
}

func (rh *AdminRouteHandler) AdminPathPrefixes() ([][]byte) {
	return rh.pathPrefixes
}

func (rh *AdminRouteHandler) HasPrefix(path []byte) (bool) {
	for _, pathPrefix := range rh.pathPrefixes {
		if bytes.HasPrefix(path, pathPrefix) {
			return true
		}
	}
	return false
}

func (rh *AdminRouteHandler) RootHandler(ctx *fasthttp.RequestCtx) {
	switch {
	case rh.baseBundle.Auth.IsLoggedOut(ctx)                  : rh.baseBundle.RouteHandler.RedirectToLogin(ctx); return
	case rh.baseBundle.Auth.DoesNotHaveAdminAuthorization(ctx): rh.baseBundle.RouteHandler.Route403(ctx); return
	}

	path   := bytes.TrimRight(ctx.Path(), "/")
	isGET  := ctx.IsGet()
	isPOST := ctx.IsPost()

	___redirectTo := rh.baseBundle.RouteHandler.RedirectTo

	switch {
	case isGET && bytes.Equal(path, []byte("/admin"))          : ___redirectTo(ctx, []byte("/admin/user/listing"))
	case isGET && bytes.Equal(path, []byte("/admin/dashboard")): rh.renderDashboard(ctx)

	case isGET && bytes.Equal(path, []byte("/admin/user"))        : ___redirectTo(ctx, []byte("/admin/user/listing"))
	case isGET && bytes.Equal(path, []byte("/admin/user/listing")): rh.renderUserListing(ctx)

	case isGET  && bytes.Equal(path, []byte("/admin/user/detail"))          : rh.renderUserDetail(ctx)
	case isPOST && bytes.Equal(path, []byte("/admin/user/detail/do/add"))   : rh.doAddUser(ctx)
	case isPOST && bytes.Equal(path, []byte("/admin/user/detail/do/edit"))  : rh.doEditUser(ctx)
	case isPOST && bytes.Equal(path, []byte("/admin/user/detail/do/delete")): rh.doDeleteUser(ctx)

	case bytes.HasPrefix(path, []byte("/admin/assets")): rh.assetHandler(ctx)

	default:
		rh.baseBundle.RouteHandler.Route404(ctx)
	}
}
