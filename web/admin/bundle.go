package admin

import (
	"noname001/logging"

	webBase "noname001/web/base"

	"noname001/web/admin/embedding"
	"noname001/web/admin/route"
)

type AdminBundleParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle
}
type AdminBundle struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	templating   *embedding.AdminTemplating
	RouteHandler *route.AdminRouteHandler
}

func NewAdminBundle(params *AdminBundleParams) (*AdminBundle, error) {
	var err error

	bundle := &AdminBundle{}
	bundle.logger = params.Logger
	bundle.logPrefix = params.LogPrefix + ".admin"

	bundle.baseBundle = params.BaseBundle

	bundle.templating, err = embedding.NewAdminTemplating(&embedding.AdminTemplatingParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,

		BaseBundle: bundle.baseBundle,
	})
	if err != nil {
		return nil, err
	}

	bundle.RouteHandler = route.NewAdminRouteHandler(&route.AdminRouteHandlerParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,

		BaseBundle: bundle.baseBundle,

		Templating: bundle.templating,
	})

	return bundle, nil
}
