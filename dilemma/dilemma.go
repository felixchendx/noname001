package dilemma

import (
	"noname001/logging"

	"noname001/dilemma/comm"
)

// containment area for:
// - dilemmatic stuffs
// - temporary stuffs
// - what-to-name stuffs
// - whatever stuffs
// - in-transit stuffs
// and yada yada yada

type LogPrefix struct {
	Prefixes  []string
	Separator string
}

type OriginTree []string // i.e. ["app", "common", "device", "coordinator", "etc"]


func InjectTemporaryLogger(_logger *logging.WrappedLogger) {
	comm.SetLogger(_logger)
}
