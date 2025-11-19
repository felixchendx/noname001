package mediasrv

import (
	"context"

	"noname001/logging"

	"noname001/app/module/common/mediasrv/filesystem"
	"noname001/app/module/common/mediasrv/config"
	"noname001/app/module/common/mediasrv/event"
	"noname001/app/module/common/mediasrv/coordinator"
)

type ModuleParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger

	Config        *config.RawModuleConfig
}
type Module struct {
	context     context.Context
	cancel      context.CancelFunc

	logger      *logging.WrappedLogger
	logPrefix   string

	config      *config.RawModuleConfig

	evHub       *event.EventHub
	coordinator *coordinator.Coordinator

	state       string
}

func NewModule(params *ModuleParams) (*Module, error) {
	var err error

	mod := &Module{}
	mod.context, mod.cancel = context.WithCancel(params.ParentContext)
	mod.logger, mod.logPrefix = params.Logger, "app.common.mediasrv"

	mod.config = params.Config

	// === event ===
	mod.evHub = event.ExtendLocalEventHub()
	// === event ===

	// === coordinator ===
	mod.coordinator, err = coordinator.NewCoordinator(&coordinator.CoordinatorParams{
		ParentContext: mod.context,
		Logger: mod.logger,
		LogPrefix: mod.logPrefix,
		Config: params.Config,
		EvHub: mod.evHub,
	})
	if err != nil {
		mod.logger.Errorf("%s: new coordinator err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === coordinator ===

	mod.state = "init"
	mod.logger.Infof("%s: initialized.", mod.logPrefix)

	return mod, nil
}

func (mod *Module) Init() (err error) {
	err = filesystem.PrepareAll()
	if err != nil { return }

	return
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
