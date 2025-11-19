package definition

type ModuleCode string
const (
	COMMON_ORG        ModuleCode = "common_org" // organizing stuffs, i.e. arbitrary grouping, etc
	COMMON_MEDIASRV   ModuleCode = "common_mediasrv"
	COMMON_CACHE      ModuleCode = "common_cache" // upgrade to 'COMMON_INTEL' or 'COMMON_INFO_BROKER'
	COMMON_DEVICE     ModuleCode = "common_device"

	FEATURE_BACKUP  ModuleCode = "feature_backup"
	FEATURE_STREAM  ModuleCode = "feature_stream"
	
	FEATURE_WALL    ModuleCode = "feature_wall"

	// FEATURE_MONITORING // TODO

	// FEATURE_AI      ModuleCode = "feature_ai" // rename to analytics ?
	FEATURE_ARCHIVE ModuleCode = "feature_archive"

	// FEATURE_CC      ModuleCode = "feature_cc"

	// PROJECT_XXX    ModuleCode = "project_xxx"
)

type ModuleState string
const (
	STATE_INIT  ModuleState = "init"
	STATE_START ModuleState = "start"
	STATE_STOP  ModuleState = "stop"
)



// unimplemented VVV
// type ModuleType string
// type ModuleDefinition struct {
// 	Code string // ALL CAPS
// 	Type ModuleType

// 	Name string // ALL LOW
// }

// const (
// 	MODTYPE_COMMON
// 	MODTYPE_FEATURE
// 	MODTYPE_PROJECT
// )
