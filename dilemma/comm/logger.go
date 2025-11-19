package comm

import (
	"noname001/logging"

	mdapi "noname001/dilemma/comm/zmqdep/mdp"
	"noname001/dilemma/comm/zmqdep/kvmsg"
)

var (
	logger *logging.WrappedLogger
)

func SetLogger(_logger *logging.WrappedLogger) {
	logger = _logger

	mdapi.SetLogger(_logger)
	kvmsg.SetLogger(_logger)
}
