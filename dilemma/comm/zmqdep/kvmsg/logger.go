package kvmsg

import (
	"noname001/logging"
)

var (
	logger *logging.WrappedLogger
)

func SetLogger(_logger *logging.WrappedLogger) {
	logger = _logger
}
