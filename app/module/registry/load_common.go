package registry

import (
	"noname001/logging"

	modDef "noname001/app/module/definition"

	modMediasrv "noname001/app/module/common/mediasrv"
	modCache    "noname001/app/module/common/cache"
	modDevice   "noname001/app/module/common/device"
)

func (modRegistry *ModuleRegistry) loadCommonMediasrv() (error) {
	modParams := &modMediasrv.ModuleParams{
		ParentContext: modRegistry.context,
		Logger: logging.Logger,

		Config: &modRegistry.cfgRoot.ModuleMediasrv,
	}

	_modMediasrv, err := modMediasrv.NewModule(modParams)
	if err != nil {
		return err
	}

	modRegistry.loadedModules[modDef.COMMON_MEDIASRV] = _modMediasrv

	return nil
}

func (modRegistry *ModuleRegistry) loadCommonCache() (error) {
	modParams := &modCache.ModuleParams{
		ParentContext: modRegistry.context,
		Logger: logging.Logger,

		// CommonParams:
		Timezone: modRegistry.cfgRoot.Global.TimeLoc,
	}

	_modCache, err := modCache.NewModule(modParams)
	if err != nil {
		return err
	}

	modRegistry.loadedModules[modDef.COMMON_CACHE] = _modCache

	return nil
}

func (modRegistry *ModuleRegistry) loadCommonDevice() (error) {
	modParams := &modDevice.ModuleParams{
		Context: modRegistry.context,
		Logger: logging.Logger,
		// LogPrefix: ""
		Config: &modRegistry.cfgRoot.ModuleDevice,

		// CommonParams: 
		Timezone: modRegistry.cfgRoot.Global.TimeLoc,
	}
	if modRegistry.runnerCfgRoot != nil {
		modParams.RunnerConfig = &modRegistry.runnerCfgRoot.ModuleDevice
	}

	_modDevice, err := modDevice.NewModule(modParams)
	if err != nil {
		return err
	}

	modRegistry.loadedModules[modDef.COMMON_DEVICE] = _modDevice

	return nil
}
