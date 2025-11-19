package registry

import (
	"noname001/logging"

	modDef "noname001/app/module/definition"

	modWall  "noname001/app/module/feature/wall"
)

// small-y modules that requires hub to functions here

func (modRegistry *ModuleRegistry) loadFeatureWall() (error) {
	modParams := &modWall.ModuleParams{
		Context: modRegistry.context,
		Logger: logging.Logger,
		Config: &modRegistry.cfgRoot.ModuleWall,

		// CommonParams: 
		Timezone: modRegistry.cfgRoot.Global.TimeLoc,
	}
	if modRegistry.runnerCfgRoot != nil {
		modParams.RunnerConfig = &modRegistry.runnerCfgRoot.ModuleWall
	}

	// DIBundle

	_modWall, err := modWall.NewModule(modParams)
	if err != nil {
		return err
	}

	modRegistry.loadedModules[modDef.FEATURE_WALL] = _modWall

	return nil
}
