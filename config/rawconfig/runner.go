package rawconfig

import (
	deviceRunner "noname001/app/module/common/device/runner"

	streamRunner "noname001/app/module/feature/stream/runner"

	wallRunner   "noname001/app/module/feature/wall/runner"
)

type RunnerConfigRoot struct {
	Something Something `yaml:"something"`

	ModuleDevice deviceRunner.RawRunnerConfig `yaml:"module_device"`

	ModuleStream streamRunner.RawRunnerConfig `yaml:"module_stream"`

	ModuleWall   wallRunner.RawRunnerConfig   `yaml:"module_wall"`
}

type Something struct {}
