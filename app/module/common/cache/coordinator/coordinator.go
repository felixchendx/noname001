package coordinator

import (
	"context"
	"time"

	"noname001/logging"

	cacheEv      "noname001/app/module/common/cache/event"
	cacheIntface "noname001/app/module/common/cache/intface"

	"noname001/app/module/common/cache/coordinator/live"
)

type CoordinatorParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger
	LogPrefix     string

	EvHub         *cacheEv.EventHub

	Timezone      *time.Location
}
type Coordinator struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	evHub     *cacheEv.EventHub

	liveCache *live.LiveCache
}

func NewCoordinator(params *CoordinatorParams) (*Coordinator, error) {
	var err error

	coord := &Coordinator{}
	coord.context, coord.cancel = context.WithCancel(params.ParentContext)
	coord.logger, coord.logPrefix = params.Logger, params.LogPrefix + ".cdt"

	coord.evHub = params.EvHub

	coord.liveCache, err = live.NewLiveCache(&live.LiveCacheParams{
		ParentContext: coord.context,
		ParentLogger : coord.logger, ParentLogPrefix: coord.logPrefix,

		EvHub: coord.evHub,

		Timezone: params.Timezone,
	})
	if err != nil {
		return nil, err
	}

	cacheIntface.AssignCacheDataProvider(coord.liveCache)

	return coord, nil
}

func (coord *Coordinator) Start() (err error) {
	err = coord.liveCache.Start()
	if err != nil { return }

	return
}

func (coord *Coordinator) PostStart() {
	coord.liveCache.PostStart()
}

func (coord *Coordinator) Stop() (err error) {
	coord.liveCache.Stop()

	return
}
