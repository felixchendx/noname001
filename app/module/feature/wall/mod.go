package wall

import (
	"context"
	"time"

	"noname001/logging"

	"noname001/app/module/feature/wall/filesystem"
	"noname001/app/module/feature/wall/config"
	"noname001/app/module/feature/wall/event"
	"noname001/app/module/feature/wall/store"
	"noname001/app/module/feature/wall/comm"
	"noname001/app/module/feature/wall/coordinator"
	"noname001/app/module/feature/wall/service"
	"noname001/app/module/feature/wall/runner"
)

type ModuleParams struct {
	Context      context.Context
	Logger       *logging.WrappedLogger
	Config       *config.RawModuleConfig
	RunnerConfig *runner.RawRunnerConfig

	Timezone     *time.Location
}
type Module struct {
	context     context.Context
	cancel      context.CancelFunc
	
	logger      *logging.WrappedLogger
	logPrefix   string

	config      *config.RawModuleConfig
	
	evHub       *event.EventHub
	store       *store.Store
	commBundle  *comm.CommBundle
	coordinator *coordinator.Coordinator
	service     *service.Service
	runner      *runner.Runner

	state       string
	// TODO: internal app evChan

	// DIBundle // TODO injected by module loader
}

func NewModule(params *ModuleParams) (*Module, error) {
	var err error
	
	mod := &Module{}
	mod.context, mod.cancel = context.WithCancel(params.Context)
	mod.logger = params.Logger
	mod.logPrefix = "app.feature.wall"

	mod.config = params.Config

	// === init ===
	err = mod.Init()
	if err != nil {
		mod.logger.Errorf("%s: init err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === init ===

	// === event ===
	mod.evHub = event.ExtendLocalEventHub()
	// === event ===

	// === store ===
	_store, storeErr := store.NewStore(&store.StoreParams{
		Context: mod.context,
		Logger: mod.logger,
		LogPrefix: mod.logPrefix,
	})
	if storeErr != nil {
		mod.logger.Errorf("%s: new store err %s", mod.logPrefix, storeErr.Error())
		mod.Stop()
		return nil, storeErr
	}
	mod.store = _store

	sev := mod.store.DB.DBInit()
	if sev.IsError() {
		mod.logger.Errorf("%s: dbinit err %s", mod.logPrefix, sev.OriErr().Error())
		mod.Stop()
		return nil, sev.OriErr()
	}

	sev = mod.store.DB.DBUpdate()
	if sev.IsError() {
		mod.logger.Errorf("%s: dbupdate err %s", mod.logPrefix, sev.OriErr().Error())
		mod.Stop()
		return nil, sev.OriErr()
	}
	// === store ===

	// === comm bundle ===
	mod.commBundle, err = comm.NewCommBundle(&comm.CommBundleParams{
		Context: mod.context,
		Logger: mod.logger,
		LogPrefix: mod.logPrefix,
		Config: params.Config,
	})
	if err != nil {
		mod.logger.Errorf("%s: new comm bundle err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === comm bundle ===

	// === coordinator ===
	mod.coordinator, err = coordinator.NewCoordinator(&coordinator.CoordinatorParams{
		Context: mod.context,
		Logger: mod.logger,
		LogPrefix: mod.logPrefix,
		Config: params.Config,
		EvHub: mod.evHub,
		Store: mod.store,
		CommBundle: mod.commBundle,
		Timezone: params.Timezone,
	})
	if err != nil {
		mod.logger.Errorf("%s: new coordinator err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === coordinator ===

	// === service ===
	mod.service, err = service.NewService(&service.ServiceParams{
		Context: mod.context,
		Logger: mod.logger,
		LogPrefix: mod.logPrefix,
		EvHub: mod.evHub,
		Store: mod.store,
		CommBundle: mod.commBundle,
		Coordinator: mod.coordinator,
	})
	if err != nil {
		mod.logger.Errorf("%s: new service err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === service ===

	// === runner ===
	mod.runner, err = runner.NewRunner(&runner.RunnerParams{
		Logger: mod.logger,
		LogPrefix: mod.logPrefix,
		Service: mod.service,
		RunnerConfig: params.RunnerConfig,
	})
	if err != nil {
		mod.logger.Errorf("%s: new runner err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === runner ===

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
	
	{
		err = mod.store.Open()
		if err != nil { return }
	}

	{
		err = mod.commBundle.Connect()
		if err != nil { return }
	}

	{
		err = mod.coordinator.Start()
		if err != nil { return }
	}

	{
		err = mod.service.Start()
		if err != nil { return }
	}

	{
		err = mod.runner.Execute()
		if err != nil { return }

		mod.runner.Destroy()
		mod.runner = nil
	}

	mod.state = "start"
	mod.logger.Infof("%s: started.", mod.logPrefix)

	return
}

func (mod *Module) PostStart() {
	mod.coordinator.PostStart()
	mod.service.PostStart()

	mod.logger.Infof("%s: post-start executed.", mod.logPrefix)
}

func (mod *Module) PreStop() {
}

func (mod *Module) Stop() (err error) {
	// stop in reverse order
	if mod.runner != nil { mod.runner.Destroy() }
	if mod.service != nil { mod.service.Stop() }
	if mod.coordinator != nil { mod.coordinator.Stop() }

	if mod.commBundle != nil { mod.commBundle.Disconnect() }
	if mod.store != nil { mod.store.Close() }
	if mod.evHub != nil { mod.evHub.Close() }

	if mod.context != nil { mod.cancel() }

	mod.state = "stop"
	mod.logger.Infof("%s: stopped.", mod.logPrefix)
	return
}

func (mod *Module) State() (string) {
	return mod.state
}
