package registry

import (
	"context"

	"noname001/config/rawconfig"

	modDef "noname001/app/module/definition"
)

type ModuleRegistry struct {
	context       context.Context
	cancel        context.CancelFunc

	// list of all available mods ordered by dependency, thus load sequence
	modList       []modDef.ModuleCode
	modDeps       map[modDef.ModuleCode][]modDef.ModuleCode
	
	// accompaniment to modList, indicating what mod will be loaded
	modToLoad     map[modDef.ModuleCode]bool
	
	// merge this to ^^^ ?
	loadedModules map[modDef.ModuleCode]modDef.ModuleIntface


	cfgRoot       *rawconfig.ConfigRoot
	runnerCfgRoot *rawconfig.RunnerConfigRoot
}

func LoadModules(ctx context.Context, cfg *rawconfig.ConfigRoot, runnerCfg *rawconfig.RunnerConfigRoot) (*ModuleRegistry, error) {
	modRegistry := &ModuleRegistry{}
	modRegistry.context, modRegistry.cancel = context.WithCancel(ctx)
	modRegistry.modList = make([]modDef.ModuleCode, 0)
	modRegistry.modDeps = make(map[modDef.ModuleCode][]modDef.ModuleCode)
	modRegistry.modToLoad = make(map[modDef.ModuleCode]bool)
	modRegistry.loadedModules = make(map[modDef.ModuleCode]modDef.ModuleIntface)

	modRegistry.cfgRoot = cfg
	modRegistry.runnerCfgRoot = runnerCfg

	modRegistry.defineModuleMasterList()
	modRegistry.defineModuleDependencies()

	modRegistry.markModulesToLoad()

	err := modRegistry.loadModules()
	if err != nil {
		modRegistry.StopAll()
		return nil, err
	}

	return modRegistry, nil
}

func (modRegistry *ModuleRegistry) StartAll() (err error) {
	for _, modCode := range modRegistry.modList {
		_mod, ok := modRegistry.loadedModules[modCode]
		if ok && _mod.State() == string(modDef.STATE_INIT) {
			err = _mod.Start()
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		modRegistry.StopAll()
		return
	}

	modRegistry.runnerCfgRoot = nil

	for _, modCode := range modRegistry.modList {
		_mod, ok := modRegistry.loadedModules[modCode]
		if ok {
			_mod.PostStart()
		}
	}

	return
}

func (modRegistry *ModuleRegistry) StopAll() (err error) {
	// prestop

	for i := len(modRegistry.modList)-1; i >= 0; i-- {
		modCode := modRegistry.modList[i]
		
		_mod, ok := modRegistry.loadedModules[modCode]
		if ok && _mod.State() == string(modDef.STATE_START) {
			err = _mod.Stop()
			if err != nil {

			} else {
				delete(modRegistry.loadedModules, modCode)
			}
		}
	} 

	return
}

func (modRegistry *ModuleRegistry) States() (map[modDef.ModuleCode]string) {
	states := make(map[modDef.ModuleCode]string)
	for modCode, mod := range modRegistry.loadedModules {
		states[modCode] = mod.State()
	}

	return states
}
func (modRegistry *ModuleRegistry) PlainStates() (map[string]string) {
	states := make(map[string]string)
	for modCode, mod := range modRegistry.loadedModules {
		states[string(modCode)] = string(mod.State())
	}

	return states
}

// TODO: move to dependencies.go, for streamline module editing
func (modRegistry *ModuleRegistry) markModulesToLoad() {
	// sweep outer modules to load
	for _, modCode := range modRegistry.modList {
		switch modCode {
		case modDef.FEATURE_STREAM:
			if modRegistry.cfgRoot.Application.RunModuleStream {
				modRegistry.modToLoad[modCode] = true
			}

		case modDef.FEATURE_WALL:
			if modRegistry.cfgRoot.Application.RunModuleWall {
				modRegistry.modToLoad[modCode] = true
			}

		default:
			modRegistry.modToLoad[modCode] = false
		}
	}

	// sweep dependencies to load
	for modCode, modToBeLoaded := range modRegistry.modToLoad {
		if modToBeLoaded {
			deps, ok := modRegistry.modDeps[modCode]
			if ok {
				for _, depCode := range deps {
					modRegistry.modToLoad[depCode] = true
				}
			}
		}
	}
}

// TODO: move to dependencies.go
func (modRegistry *ModuleRegistry) loadModules() (err error) {
	for _, modCode := range modRegistry.modList {
		modToBeLoaded, ok := modRegistry.modToLoad[modCode]
		if ok && modToBeLoaded {
			switch modCode {
			case modDef.COMMON_CACHE   : err = modRegistry.loadCommonCache()
			case modDef.COMMON_MEDIASRV: err = modRegistry.loadCommonMediasrv()
			case modDef.COMMON_DEVICE  : err = modRegistry.loadCommonDevice()

			case modDef.FEATURE_STREAM: err = modRegistry.loadFeatureStream()

			case modDef.FEATURE_WALL:   err = modRegistry.loadFeatureWall()
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
