package live

import (
	"time"
)

func (lc *LiveCache) newMediaServer() (*t_mediasrv) {
	return &t_mediasrv{
		ipToUse  : "",
		ports    : make(map[string]string),
		authnPair: "",

		// pingTs: zeroTime,
		pingResults: make([]*t_pingResult, 0),

		opFlag  : false,
		opStatus: "",
	}
}

// -----------------------------------------------------------------------------
type t_mediasrv struct {
	ipToUse   string
	ports     map[string]string
	authnPair string

	pingTs      time.Time
	pingResults []*t_pingResult

	opFlag   bool
	opStatus string
}

type t_pingResult struct {
	ip        string
	reachable bool
	err       error
}
