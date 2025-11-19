package runner

import (
	"noname001/logging"

	"noname001/app/module/common/device/service"
)

type RunnerParams struct {
	Logger       *logging.WrappedLogger
	LogPrefix    string
	Service      *service.Service
	RunnerConfig *RawRunnerConfig
}
type Runner struct {
	logger       *logging.WrappedLogger
	logPrefix    string
	service      *service.Service
	runnerConfig *RawRunnerConfig
}

func NewRunner(params *RunnerParams) (*Runner, error) {
	runna := &Runner{}
	runna.logger = params.Logger
	runna.logPrefix = params.LogPrefix + ".runner"
	runna.service = params.Service
	runna.runnerConfig = params.RunnerConfig

	return runna, nil
}

func (runna *Runner) Execute() (err error) {
	err = runna.exec()
	if err != nil { return }

	return
}

func (runna *Runner) Destroy() (err error) {
	return nil
}
