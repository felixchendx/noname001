package device

import (
	"context"
	"time"

	"noname001/logging"

	"noname001/app/base/sec"
	"noname001/app/module/common/device/config"
	"noname001/app/module/common/device/event"
	"noname001/app/module/common/device/store"
	"noname001/app/module/common/device/comm"
	"noname001/app/module/common/device/coordinator"
	"noname001/app/module/common/device/service"
	"noname001/app/module/common/device/runner"

	"noname001/app/module/common/device/intface"
)

type ModuleParams struct {
	Context      context.Context
	Logger       *logging.WrappedLogger
	// LogPrefix    dilemma.LogPrefix
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

	secBundle   *sec.DumbSecurityBundle

	evHub       *event.EventHub
	store       *store.Store
	commBundle  *comm.CommBundle
	coordinator *coordinator.Coordinator
	service     *service.Service
	runner      *runner.Runner

	state       string
}

func NewModule(params *ModuleParams) (*Module, error) {
	var err error
	
	mod := &Module{}
	mod.context, mod.cancel = context.WithCancel(params.Context)
	mod.logger = params.Logger
	mod.logPrefix = "app.common.device"

	mod.config = params.Config

	mod.secBundle = sec.NewDumbSecurityBundle()
	// TODO: move to config when provz is kinda established
	mod.secBundle.AddKey("device", "f7cdc572cc956e7f0f590c8894d82d56b64d2c44b7aa05326237cbb14ab18e0b")

	// === init ===
	err = mod.Init()
	if err != nil {
		mod.logger.Errorf("%s: init err %s", mod.logPrefix, err.Error())
		mod.Stop()
		return nil, err
	}
	// === init ===

	// === event ===
	mod.evHub = event.NewEventHub(&event.EventHubParams{
		ParentContext: mod.context,
		Logger: mod.logger, LogPrefix: mod.logPrefix,
	})
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
		SecBundle: mod.secBundle,
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
		SecBundle: mod.secBundle,
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

	intface.AssignDeviceEventProvider(mod.evHub)
	mod.commBundle.DeviceProvider.AssignDataHandler(mod.service)

	mod.state = "init"
	mod.logger.Infof("%s: initialized.", mod.logPrefix)
	return mod, nil
}

func (mod *Module) Init() (err error) {
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
