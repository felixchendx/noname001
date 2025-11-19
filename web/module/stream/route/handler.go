package route

import (
	"bytes"

	"github.com/valyala/fasthttp"

	"noname001/logging"
	webBase "noname001/web/base"
	"noname001/web/module/stream/embedding"
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
	var rh *ModuleRouteHandler = &ModuleRouteHandler{}
	rh.logger = params.Logger
	rh.logPrefix = params.LogPrefix + ".rh"

	rh.baseBundle = params.BaseBundle

	rh.templating = params.Templating

	rh.pathPrefixes = [][]byte{
		[]byte("/stream"),
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
	case rh.baseBundle.Navi.DoesNotHave("stream"): rh.baseBundle.RouteHandler.Route404_ModuleNotActive(ctx, "stream"); return
	case rh.baseBundle.Auth.IsLoggedOut(ctx)     : rh.baseBundle.RouteHandler.RedirectToLogin(ctx); return
	}

	path := bytes.TrimRight(ctx.Path(), "/")
	isGET := ctx.IsGet()
	isPOST := ctx.IsPost()

	___redirectTo := rh.baseBundle.RouteHandler.RedirectTo

	switch {
	case isGET && bytes.Equal(path, []byte("/stream"))                       : ___redirectTo(ctx, []byte("/stream/stream-group/listing"))
	case isGET && bytes.Equal(path, []byte("/stream/dashboard"))             : rh.renderDashboard(ctx)
	
	case isGET && bytes.Equal(path, []byte("/stream/stream-profile"))        : ___redirectTo(ctx, []byte("/stream/stream-profile/listing"))
	case isGET && bytes.Equal(path, []byte("/stream/stream-profile/listing")): rh.renderStreamProfileListing(ctx)

	case isGET && bytes.Equal(path, []byte("/stream/stream-profile/detail"))           : rh.renderStreamProfileDetail(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-profile/detail/do/add"))   : rh.doAddStreamProfile(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-profile/detail/do/edit"))  : rh.doEditStreamProfile(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-profile/detail/do/delete")): rh.doDeleteStreamProfile(ctx)


	case isGET && bytes.Equal(path, []byte("/stream/stream-group"))        : ___redirectTo(ctx, []byte("/stream/stream-group/listing"))
	case isGET && bytes.Equal(path, []byte("/stream/stream-group/listing")): rh.renderStreamGroupListing(ctx)

	case isGET && bytes.Equal(path, []byte("/stream/stream-group/detail"))           : rh.renderStreamGroupDetail(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-group/detail/do/add"))   : rh.doAddStreamGroup(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-group/detail/do/edit"))  : rh.doEditStreamGroup(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-group/detail/do/delete")): rh.doDeleteStreamGroup(ctx)

	case isPOST && bytes.Equal(path, []byte("/stream/local-api/add-stream-item"))   : rh.localAPI_addStreamItem(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/local-api/edit-stream-item"))  : rh.localAPI_editStreamItem(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/local-api/delete-stream-item")): rh.localAPI_deleteStreamItem(ctx)

	case isPOST && bytes.Equal(path, []byte("/stream/local-api/device-channel-preview")): rh.localAPI__requestDeviceChannelPreview(ctx)

	case isGET && bytes.Equal(path, []byte("/stream/stream-group/detail-02"))           : rh.streamGropuDetail__render(ctx)
	case isGET  && bytes.Equal(path, []byte("/stream/stream-group/detail-02/ws"))       : rh.streamGroupDetail__ws(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-group/detail-02/do/add"))   : rh.streamGroupDetail__doAdd(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-group/detail-02/do/edit"))  : rh.streamGroupDetail__doEdit(ctx)
	case isPOST && bytes.Equal(path, []byte("/stream/stream-group/detail-02/do/delete")): rh.streamGroupDetail__doDelete(ctx)

	case bytes.HasPrefix(path, []byte("/stream/assets")): rh.assetHandler(ctx)

	default:
		rh.baseBundle.RouteHandler.Route404(ctx)
	}
}
