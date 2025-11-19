package monitoring

import (
	"noname001/logging"
	webBase "noname001/web/base"
	"noname001/web/module/monitoring/embedding"
	"noname001/web/module/monitoring/route"
)

type ModuleBundleParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle
}

type ModuleBundle struct {
	logger       *logging.WrappedLogger
	logPrefix    string
	
	baseBundle   *webBase.BaseBundle
	templating   *embedding.ModuleTemplating
	RouteHandler *route.ModuleRouteHandler
}

func NewModuleBundle(params *ModuleBundleParams) (*ModuleBundle, error) {
	var err error

	bundle := &ModuleBundle{}
	bundle.logger = params.Logger
	bundle.logPrefix = params.LogPrefix + ".monit"

	bundle.baseBundle = params.BaseBundle

	bundle.templating, err = embedding.NewModuleTemplating(&embedding.ModuleTemplatingParams{
		Logger:    bundle.logger,
		LogPrefix: bundle.logPrefix,

		BaseBundle: bundle.baseBundle,
	})
	if err != nil {
		return nil, err
	}

	bundle.RouteHandler = route.NewModuleRouteHandler(&route.ModuleRouteHandlerParams{
		Logger:    bundle.logger,
		LogPrefix: bundle.logPrefix,

		BaseBundle: bundle.baseBundle,

		Templating: bundle.templating,
	})

	return bundle, nil
}
