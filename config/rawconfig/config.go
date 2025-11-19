package rawconfig

import (
	"time"

	modMediasrvConfig "noname001/app/module/common/mediasrv/config"
	modDeviceConfig   "noname001/app/module/common/device/config"

	modStreamConfig "noname001/app/module/feature/stream/config"

	modWallConfig   "noname001/app/module/feature/wall/config"
)

type ConfigRoot struct {
	ModeMock bool

	CfgDef      ConfigDefinition `ini:"configuration_definition"`
	Logging     Logging          `ini:"logging"`
	Global      Global           `ini:"global"`

	Node        Node             `ini:"node"`
	Hub         Hub              `ini:"hub"`
	Application Application      `ini:"application"`
	Web         RawWebConfig     `ini:"web"`

	ModuleMediasrv modMediasrvConfig.RawModuleConfig `ini:"module_mediasrv"`
	ModuleDevice   modDeviceConfig.RawModuleConfig `ini:"module_device"`

	ModuleStream modStreamConfig.RawModuleConfig `ini:"module_stream"`
	ModuleWall   modWallConfig.RawModuleConfig   `ini:"module_wall"`
}

type ConfigDefinition struct {
	Version string `ini:"version"`
}

type Logging struct {
	LogTo       string `ini:"log_to"`
	LogLevel    string `ini:"log_level"`
	LogTimezone string `ini:"log_timezone"`
}

type Global struct {
	RootDirectory string `ini:"root_directory"`

	Timezone      string `ini:"timezone"`

	// injected
	TimeLoc       *time.Location
}

// TODO: restruct from end user's perspective
type Node struct {
	ID                  string `ini:"id"`
	Name                string `ini:"name"`

	Standalone          bool   `ini:"standalone"`

	BrokerServerHost    string `ini:"broker_server_host"`
	SnapshotServerHost  string `ini:"snapshot_server_host"`
	PublisherServerHost string `ini:"publisher_server_host"`
	CollectorServerHost string `ini:"collector_server_host"`
	CommVerbose         bool   `ini:"comm_verbose"`
}

type Hub struct {
	Enabled                  bool   `ini:"enabled"`

	BrokerHost               string `ini:"broker_host"`

	SnapshotHost             string `ini:"snapshot_host"`
	PublisherHost            string `ini:"publisher_host"`
	CollectorHost            string `ini:"collector_host"`

	CommVerbose              bool   `ini:"comm_verbose"`
	TmpHubProviderBrokerHost string `ini:"tmphub_provider_broker_host"`
}

type Application struct {
	Enabled         bool `ini:"enabled"`

	RunModuleBackup bool `ini:"run_module_backup"`
	RunModuleStream bool `ini:"run_module_stream"`
	RunModuleWall   bool `ini:"run_module_wall"`
}

type RawWebConfig struct {
	Enabled  bool   `ini:"enabled"`

	Hostname string `ini:"host"`
	Port     string `ini:"port"`

	Behind7Proxies bool `ini:"behind7proxies"`
}
