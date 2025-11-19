package coordinator

import (
	"context"
	"time"

	"noname001/logging"
	
	"noname001/app/module/feature/wall/config"
	"noname001/app/module/feature/wall/event"
	"noname001/app/module/feature/wall/store"
	"noname001/app/module/feature/wall/comm"
)

type CoordinatorParams struct {
	Context    context.Context
	Logger     *logging.WrappedLogger
	LogPrefix  string
	Config     *config.RawModuleConfig
	EvHub      *event.EventHub
	Store      *store.Store
	CommBundle *comm.CommBundle
	Timezone   *time.Location
}
type Coordinator struct {
	context    context.Context
	cancel     context.CancelFunc

	logger     *logging.WrappedLogger
	logPrefix  string
	
	evHub      *event.EventHub
	store      *store.Store
	commBundle *comm.CommBundle

	timezone   *time.Location
}

func NewCoordinator(params *CoordinatorParams) (*Coordinator, error) {
	var err error

	coord := &Coordinator{}
	coord.context, coord.cancel = context.WithCancel(params.Context)
	coord.logger = params.Logger
	coord.logPrefix = params.LogPrefix + ".cdt"

	coord.evHub = params.EvHub
	coord.store = params.Store
	coord.commBundle = params.CommBundle

	coord.timezone = params.Timezone

	if err != nil {
		return nil, err
	}

	return coord, nil
}

func (coord *Coordinator) Start() (err error) {
	return
}

func (coord *Coordinator) PostStart() {
}

func (coord *Coordinator) Stop() (err error) {
	// TODO: recheck all other coordinator stop routines

	coord.cancel()
	return
}
