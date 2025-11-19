package user

import (
	"noname001/logging"

	webBase "noname001/web/base"

	"noname001/web/user/embedding"
	"noname001/web/user/route"
)

type UserBundleParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle
}
type UserBundle struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	templating   *embedding.UserTemplating
	RouteHandler *route.UserRouteHandler
}

func NewUserBundle(params *UserBundleParams) (*UserBundle, error) {
	var err error

	bundle := &UserBundle{}
	bundle.logger = params.Logger
	bundle.logPrefix = params.LogPrefix + ".user"

	bundle.baseBundle = params.BaseBundle

	bundle.templating, err = embedding.NewUserTemplating(&embedding.UserTemplatingParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,

		BaseBundle: bundle.baseBundle,
	})
	if err != nil {
		return nil, err
	}

	bundle.RouteHandler = route.NewUserRouteHandler(&route.UserRouteHandlerParams{
		Logger: bundle.logger,
		LogPrefix: bundle.logPrefix,

		BaseBundle: bundle.baseBundle,

		Templating: bundle.templating,
	})

	return bundle, nil
}
