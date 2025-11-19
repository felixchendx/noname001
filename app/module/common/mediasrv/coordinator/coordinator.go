package coordinator

import (
	"context"

	"noname001/logging"

	"noname001/app/module/common/mediasrv/config"
	"noname001/app/module/common/mediasrv/event"

	mediasrvIntface "noname001/app/module/common/mediasrv/intface"

	"noname001/app/module/common/mediasrv/coordinator/mediamtxserver"
)

type CoordinatorParams struct {
	ParentContext context.Context

	Logger        *logging.WrappedLogger
	LogPrefix     string

	Config        *config.RawModuleConfig

	EvHub         *event.EventHub
}
type Coordinator struct {
	context     context.Context
	cancel      context.CancelFunc

	logger      *logging.WrappedLogger
	logPrefix   string

	evHub       *event.EventHub

	mediaServer *mediamtxserver.MediaMTXServer
}

func NewCoordinator(params *CoordinatorParams) (*Coordinator, error) {
	var err error

	coord := &Coordinator{}
	coord.context, coord.cancel = context.WithCancel(params.ParentContext)
	coord.logger, coord.logPrefix = params.Logger, params.LogPrefix + ".cdt"

	coord.evHub = params.EvHub

	coord.mediaServer, err = mediamtxserver.NewMediaMTXServer(&mediamtxserver.MediaMTXServerParams{
		Context: coord.context,
		Logger: coord.logger,
		LogPrefix: coord.logPrefix,
		Config: params.Config,
	})
	if err != nil {
		return nil, err
	}

	err = coord.mediaServer.UnpackMediamtx()
	if err != nil {
		return nil, err
	}

	mediasrvIntface.AssignMediasrvProvider(coord)

	return coord, nil
}

func (coord *Coordinator) Start() (err error) {
	err = coord.mediaServer.Start()
	if err != nil {
		return
	}

	return
}

func (coord *Coordinator) PostStart() {}

func (coord *Coordinator) Stop() (err error) {
	coord.cancel()
	return
}
