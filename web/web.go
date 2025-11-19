package web

import (
	"context"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"

	"noname001/logging"
	"noname001/version"

	webBase "noname001/web/base"

	"noname001/web/admin"
	"noname001/web/module/device"
	"noname001/web/module/stream"
	"noname001/web/module/wall"
	"noname001/web/user"

	"noname001/web/module/monitoring"
)

// TODO: https://github.com/caddyserver/caddy

type WebServiceParams struct {
	RootContext   context.Context
	RootLogger    *logging.WrappedLogger
	RootLogPrefix string

	Name       string
	ListenHost string

	Behind7Proxies  bool
	ServerStartTime time.Time

	TempModStates map[string]string
}
type WebService struct {
	name      string
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	listenHost      string
	behind7Proxies  bool
	serverStartTime time.Time

	httpServer *fasthttp.Server

	baseBundle   *webBase.BaseBundle
	adminBundle  *admin.AdminBundle
	userBundle   *user.UserBundle
	deviceBundle *device.ModuleBundle
	streamBundle *stream.ModuleBundle
	wallBundle   *wall.ModuleBundle

	monitoringBundle *monitoring.ModuleBundle
}

// TODO: do not load all bundle if the module is not active
//       modular load + service injection to prevent wild access

func NewWebService(params *WebServiceParams) (*WebService, error) {
	var err error

	svc := &WebService{}
	svc.name = params.Name
	svc.context, svc.cancel = context.WithCancel(params.RootContext)
	svc.logger, svc.logPrefix = params.RootLogger, fmt.Sprintf("%s.web", params.RootLogPrefix)

	svc.listenHost = params.ListenHost
	svc.behind7Proxies = params.Behind7Proxies
	svc.serverStartTime = params.ServerStartTime

	svc.httpServer = &fasthttp.Server{
		Logger:  svc.logger,
		Handler: svc.requestHandler,
		// ErrorHandler: func(ctx *RequestCtx, err error)
		// ConnState func(net.Conn, ConnState)
		// FormValueFunc FormValueFunc
		Name: svc.name,

		// TODO: moar configs
		// https://pkg.go.dev/github.com/valyala/fasthttp#Server
	}

	svc.baseBundle, err = webBase.NewBaseBundle(&webBase.BaseBundleParams{
		ParentContext: svc.context,
		Logger:        svc.logger,
		LogPrefix:     svc.logPrefix,

		Behind7Proxies:  svc.behind7Proxies,
		ServerStartTime: svc.serverStartTime,
		AppVersion:      version.FullVersion(),

		TempModStates: params.TempModStates,
	})
	if err != nil {
		svc.logger.Errorf("[%s] webBase.NewBaseBundle err, %s", svc.logPrefix, err.Error())
		return nil, err
	}

	svc.adminBundle, err = admin.NewAdminBundle(&admin.AdminBundleParams{
		Logger:    svc.logger,
		LogPrefix: svc.logPrefix,

		BaseBundle: svc.baseBundle,
	})
	if err != nil {
		svc.logger.Errorf("[%s] admin.NewAdminBundle err, %s", svc.logPrefix, err.Error())
		return nil, err
	}

	svc.userBundle, err = user.NewUserBundle(&user.UserBundleParams{
		Logger:    svc.logger,
		LogPrefix: svc.logPrefix,

		BaseBundle: svc.baseBundle,
	})
	if err != nil {
		svc.logger.Errorf("[%s] user.NewUserBundle err, %s", svc.logPrefix, err.Error())
		return nil, err
	}

	svc.deviceBundle, err = device.NewModuleBundle(&device.ModuleBundleParams{
		Logger:    svc.logger,
		LogPrefix: svc.logPrefix,

		BaseBundle: svc.baseBundle,
	})
	if err != nil {
		svc.logger.Errorf("[%s] device.NewModuleBundle err, %s", svc.logPrefix, err.Error())
		return nil, err
	}

	svc.streamBundle, err = stream.NewModuleBundle(&stream.ModuleBundleParams{
		Logger:     svc.logger,
		LogPrefix:  svc.logPrefix,
		BaseBundle: svc.baseBundle,
	})
	if err != nil {
		svc.logger.Errorf("[%s] stream.NewModuleBundle err, %s", svc.logPrefix, err.Error())
		return nil, err
	}

	svc.wallBundle, err = wall.NewModuleBundle(&wall.ModuleBundleParams{
		Logger:    svc.logger,
		LogPrefix: svc.logPrefix,

		BaseBundle: svc.baseBundle,
	})
	if err != nil {
		svc.logger.Errorf("[%s] wall.NewModuleBundle err, %s", svc.logPrefix, err.Error())
		return nil, err
	}

	svc.monitoringBundle, err = monitoring.NewModuleBundle(&monitoring.ModuleBundleParams{
		Logger:    svc.logger,
		LogPrefix: svc.logPrefix,

		BaseBundle: svc.baseBundle,
	})
	if err != nil {
		svc.logger.Errorf("[%s] monitoring.NewModuleBundle err, %s", svc.logPrefix, err.Error())
		return nil, err
	}

	svc.logger.Infof("[%s] initialized", svc.logPrefix)
	return svc, nil
}

func (svc *WebService) Serve() {
	go func() {
		err := svc.httpServer.ListenAndServe(svc.listenHost)
		if err != nil {
			svc.cancel()

			svc.logger.Errorf("[%s] httpserver listen err, %s", svc.logPrefix, err.Error())
			return
		}
	}()

	svc.logger.Infof("[%s] started, listening on %s", svc.logPrefix, svc.listenHost)
}

func (svc *WebService) Shutdown() {
	svc.cancel()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: issue websocket closing for all active connection
	if err := svc.httpServer.ShutdownWithContext(ctx); err != nil {
		svc.logger.Warnf("[%s] forced to shutdown: %s", svc.logPrefix, err.Error())
	}

	svc.logger.Infof("[%s] stopped", svc.logPrefix)
}

func (svc *WebService) InjectTempModStates(modStates map[string]string) {
	svc.baseBundle.InjectTempModStates(modStates)
}
