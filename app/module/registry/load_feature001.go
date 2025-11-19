package registry

import (
	"noname001/logging"

	modDef "noname001/app/module/definition"

	modStream "noname001/app/module/feature/stream"
)

func (modRegistry *ModuleRegistry) loadFeatureStream() (error) {
	modParams := &modStream.ModuleParams{
		Context: modRegistry.context,
		Logger: logging.Logger,
		Config: &modRegistry.cfgRoot.ModuleStream,

		// CommonParams: 
		Timezone: modRegistry.cfgRoot.Global.TimeLoc,
	}
	if modRegistry.runnerCfgRoot != nil {
		modParams.RunnerConfig = &modRegistry.runnerCfgRoot.ModuleStream
	}

	// DIBundle

	_modStream, err := modStream.NewModule(modParams)
	if err != nil {
		return err
	}

	modRegistry.loadedModules[modDef.FEATURE_STREAM] = _modStream

	return nil
}
