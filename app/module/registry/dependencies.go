package registry

import (
	modDef "noname001/app/module/definition"
)

func (modRegistry *ModuleRegistry) defineModuleMasterList() {
	modRegistry.modList = []modDef.ModuleCode{
		modDef.COMMON_MEDIASRV,
		modDef.COMMON_CACHE,
		modDef.COMMON_DEVICE,

		modDef.FEATURE_STREAM,
		modDef.FEATURE_BACKUP,

		modDef.FEATURE_WALL,
	}
}

func (modRegistry *ModuleRegistry) defineModuleDependencies() {
	modRegistry.modDeps[modDef.COMMON_DEVICE] = []modDef.ModuleCode{
		modDef.COMMON_MEDIASRV,
		// modDef.COMMON_CACHE,
	}

	modRegistry.modDeps[modDef.FEATURE_BACKUP] = []modDef.ModuleCode{
		modDef.COMMON_DEVICE,
	}

	modRegistry.modDeps[modDef.FEATURE_STREAM] = []modDef.ModuleCode{
		modDef.COMMON_MEDIASRV,
		// modDef.COMMON_CACHE,
		modDef.COMMON_DEVICE,
	}

	modRegistry.modDeps[modDef.FEATURE_WALL] = []modDef.ModuleCode{
		modDef.COMMON_MEDIASRV,
		modDef.COMMON_CACHE,
	}
}
