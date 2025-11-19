package route

import (
	"bytes"

	"github.com/valyala/fasthttp"
	
	"noname001/logging"

	"noname001/web/base/cookie"
	"noname001/web/base/flash"
	"noname001/web/base/auth"
	"noname001/web/base/embedding"
)

type BaseRouteHandlerParams struct {
	Logger       *logging.WrappedLogger
	LogPrefix    string

	CookieStore  *cookie.CookieStore
	FlashStore   *flash.FlashStore
	AuthProvider *auth.AuthProvider
	Templating   *embedding.BaseTemplating
}
type BaseRouteHandler struct {
	logger        *logging.WrappedLogger
	logPrefix     string

	cookieStore   *cookie.CookieStore
	flashStore    *flash.FlashStore
	authProvider  *auth.AuthProvider
	templating    *embedding.BaseTemplating

	pathPrefixes  [][]byte

	assetHandler  fasthttp.RequestHandler
}

func NewBaseRouteHandler(params *BaseRouteHandlerParams) (*BaseRouteHandler) {
	rh := &BaseRouteHandler{}
	rh.logger = params.Logger
	rh.logPrefix = params.LogPrefix + ".rh"

	rh.cookieStore = params.CookieStore
	rh.flashStore = params.FlashStore
	rh.authProvider = params.AuthProvider
	rh.templating = params.Templating

	rh.pathPrefixes = [][]byte{

	}

	rh.assetHandler = embedding.FakeSingletonBaseAssetHandler()

	return rh
}

func (rh *BaseRouteHandler) PathPrefixes() ([][]byte) {
	return rh.pathPrefixes
}

func (rh *BaseRouteHandler) HasPrefix(path []byte) (bool) {
	for _, pathPrefix := range rh.pathPrefixes {
		if bytes.HasPrefix(path, pathPrefix) {
			return true
		}
	}
	return false
}

func (rh *BaseRouteHandler) RouteAsset(ctx *fasthttp.RequestCtx) {
	rh.assetHandler(ctx)
}
