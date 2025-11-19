package coordinator

import (
	"context"
	"time"

	"noname001/logging"
	
	"noname001/app/module/feature/stream/config"
	"noname001/app/module/feature/stream/event"
	"noname001/app/module/feature/stream/store"
	"noname001/app/module/feature/stream/comm"

	"noname001/app/module/feature/stream/coordinator/live"
	"noname001/app/module/feature/stream/coordinator/preview"
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
	context     context.Context
	cancel      context.CancelFunc

	logger      *logging.WrappedLogger
	logPrefix   string

	evHub       *event.EventHub
	store       *store.Store
	commBundle  *comm.CommBundle

	timezone    *time.Location

	liveStreams map[string]*live.LiveStream // need mutex ?

	dyingPreviews map[string]*preview.DeviceChannelPreview
}

func NewCoordinator(params *CoordinatorParams) (*Coordinator, error) {
	// var err error

	coord := &Coordinator{}
	coord.context, coord.cancel = context.WithCancel(params.Context)
	coord.logger = params.Logger
	coord.logPrefix = params.LogPrefix + ".cdt"

	coord.evHub = params.EvHub
	coord.store = params.Store
	coord.commBundle = params.CommBundle

	coord.timezone = params.Timezone

	coord.liveStreams = make(map[string]*live.LiveStream)

	coord.dyingPreviews =  make(map[string]*preview.DeviceChannelPreview)

	return coord, nil
}

func (coord *Coordinator) Start() (err error) {
	coord.dcpCleanupWorker()

	return
}

func (coord *Coordinator) PostStart() {
	coord.postStartRoutine()

	go coord.streamEventListeners()
}

func (coord *Coordinator) Stop() (err error) {
	coord.cancel()
	return
}
