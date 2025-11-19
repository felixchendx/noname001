package live

import (
	"time"
)

// is a blocking call
func (lc *LiveCache) startupRoutine() {
	// timeout here, or modify mod/registry for stuffs that need to be started in specific orders
	// this artificial timeout is necessary for 
	// nodes that run all mod in 1 node (i.e. local dev or mini setup)
	time.Sleep(1 * time.Second)

	// populate initial data and wait for it to complete before starting other workers
	// to prevent jumbled / unnecessary network calls
	lc.scanRoutine()

	go lc.liveCacheWorker()
	go lc.crossnodeEventListeners()

	lc.setupCrons()
}
