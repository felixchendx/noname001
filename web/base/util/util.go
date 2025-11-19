package util

import (
	"noname001/logging"
)

type UtilParams struct {
	Logger    *logging.WrappedLogger
	LogPrefix string
}
type Util struct {
	logger    *logging.WrappedLogger
	logPrefix string
}

func NewUtil(params *UtilParams) (*Util) {
	util := &Util{}
	util.logger = params.Logger
	util.logPrefix = params.LogPrefix + ".util"

	return util
}
