package cache

import (
	"context"
	"time"

	"noname001/logging"

	"noname001/app/module/common/cache/event"
	"noname001/app/module/common/cache/coordinator"

	"noname001/app/module/common/cache/intface"
)

type ModuleParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger

	Timezone      *time.Location
}
type Module struct {
	context     context.Context
	cancel      context.CancelFunc

	logger      *logging.WrappedLogger
	logPrefix   string

	evHub       *event.EventHub
	coordinator *coordinator.Coordinator

	state       string
}

func NewModule(params *ModuleParams) (*Module, error) {
	var err error

	mod := &Module{}
	mod.context, mod.cancel = context.WithCancel(params.ParentContext)
	mod.logger, mod.logPrefix = params.Logger, "app.common.cache"

	// === event ===
	mod.evHub = event.NewEventHub(&event.EventHubParams{
		ParentContext: mod.context,
		Logger       : mod.logger, LogPrefix: mod.logPrefix,
	})
	// === event ===

	// === coordinator ===
	mod.coordinator, err = coordinator.NewCoordinator(&coordinator.CoordinatorParams{
		ParentContext: mod.context,
		Logger: mod.logger,
		LogPrefix: mod.logPrefix,
		EvHub: mod.evHub,

		Timezone: params.Timezone,
	})
	if err != nil {
		mod.logger.Errorf("%s: new coordinator err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === coordinator ===

	mod.state = "init"
	mod.logger.Infof("%s: initialized.", mod.logPrefix)

	intface.AssignCacheEventProvider(mod.evHub)

	return mod, nil
}

func (mod *Module) Start() (err error) {
	mod.evHub.Open()

	err = mod.coordinator.Start()
	if err != nil { return }

	mod.state = "start"
	mod.logger.Infof("%s: started.", mod.logPrefix)

	return
}

func (mod *Module) PostStart() {
	mod.coordinator.PostStart()

	mod.logger.Infof("%s: post-start executed.", mod.logPrefix)
}

func (mod *Module) PreStop() {}

func (mod *Module) Stop() (err error) {
	// stop in reverse order
	if mod.coordinator != nil { mod.coordinator.Stop() }

	if mod.evHub != nil { mod.evHub.Close() }
	
	if mod.context != nil { mod.cancel() }

	mod.state = "stop"
	mod.logger.Infof("%s: stopped.", mod.logPrefix)

	return
}

func (mod *Module) State() (string) {
	return mod.state
}
