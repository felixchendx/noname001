package base

import (
	"context"
	"time"
	
	"noname001/logging"

	"noname001/web/base/cookie"
	"noname001/web/base/flash"
	"noname001/web/base/auth"
	"noname001/web/base/navigation"
	"noname001/web/base/embedding"
	"noname001/web/base/route"
	"noname001/web/base/util"
	"noname001/web/base/comm/ws"
)

type BaseBundleParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger
	LogPrefix     string

	Behind7Proxies  bool
	ServerStartTime time.Time
	AppVersion      string

	TempModStates map[string]string
}
type BaseBundle struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	behind7Proxies  bool
	serverStartTime time.Time
	appVersion      string
	
	Cookie       *cookie.CookieStore
	Flash        *flash.FlashStore
	Auth         *auth.AuthProvider
	Navi         *navigation.BaseNavi

	Templating   *embedding.BaseTemplating
	RouteHandler *route.BaseRouteHandler

	Util         *util.Util

	WSHub        *ws.WSHub
}

func NewBaseBundle(params *BaseBundleParams) (*BaseBundle, error) {
	var err error

	bundle := &BaseBundle{}
	bundle.logger = params.Logger
	bundle.logPrefix = params.LogPrefix + ".base"

	bundle.behind7Proxies = params.Behind7Proxies
	bundle.serverStartTime = params.ServerStartTime
	bundle.appVersion = params.AppVersion

	bundle.Cookie = cookie.NewCookieStore(&cookie.CookieStoreParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,
	})

	bundle.Flash = flash.NewFlashStore(&flash.FlashStoreParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,
		CookieStore: bundle.Cookie,
	})

	bundle.Auth = auth.NewAuthProvider(&auth.AuthProviderParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,
		CookieStore: bundle.Cookie,
		FlashStore: bundle.Flash,
	})

	bundle.Navi = navigation.NewBaseNavi(&navigation.BaseNaviParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,

		AuthProvider: bundle.Auth,
		TempModStates: params.TempModStates,
	})

	bundle.Templating, err = embedding.NewBaseTemplating(&embedding.BaseTemplatingParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,

		AuthProvider: bundle.Auth,
		Navi: bundle.Navi,

		CacheBusterVar: params.ServerStartTime.Format("20060102T150405Z"),
		AppVersion: bundle.appVersion,
	})

	if err != nil {
		return nil, err
	}

	bundle.RouteHandler = route.NewBaseRouteHandler(&route.BaseRouteHandlerParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,
		CookieStore: bundle.Cookie,
		FlashStore: bundle.Flash,
		AuthProvider: bundle.Auth,
		Templating: bundle.Templating,
	})

	bundle.Util = util.NewUtil(&util.UtilParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,
	})

	bundle.WSHub = ws.NewWSHub(&ws.WSHubParams{
		ParentContext: params.ParentContext,
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,
	})

	return bundle, nil
}

func (baseBundle *BaseBundle) Behind7Proxies() (bool) {
	return baseBundle.behind7Proxies
}

func (baseBundle *BaseBundle) InjectTempModStates(modStates map[string]string) {
	baseBundle.Navi.InjectTempModStates(modStates)
}
