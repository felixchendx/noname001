package comm

import (
	"context"

	"noname001/logging"

	"noname001/app/module/feature/wall/config"
)

type CommBundleParams struct {
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
	Config    *config.RawModuleConfig
}

type CommBundle struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string
}

func NewCommBundle(params *CommBundleParams) (*CommBundle, error) {
	commBundle := &CommBundle{}
	commBundle.context, commBundle.cancel = context.WithCancel(params.Context)
	commBundle.logger = params.Logger
	commBundle.logPrefix = params.LogPrefix + ".comm"

	return commBundle, nil
}

func (commBundle *CommBundle) Connect() (err error) {
	return
}

func (commBundle *CommBundle) Disconnect() (err error) {
	return
}
