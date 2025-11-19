package live

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"noname001/logging"

	cacheEv "noname001/app/module/common/cache/event"
)

// when necessary, replace these internal + external data holder with inmem sqlite
// eliminating a lot of concurrency handlings (+ better perfs ?)
// but... way way more codes (minimal requirement... the whole db/store layer)

type LiveCacheParams struct {
	ParentContext   context.Context
	ParentLogger    *logging.WrappedLogger
	ParentLogPrefix string

	EvHub           *cacheEv.EventHub

	Timezone        *time.Location
}

type LiveCache struct {
	context    context.Context
	cancel     context.CancelFunc
	logger     *logging.WrappedLogger
	logPrefix  string

	evHub      *cacheEv.EventHub

	cron     *cron.Cron
	cronJobs map[string]cron.EntryID

	commBundle *t_commBundle

	miscJobChan chan *t_miscJob

	// === nodesMutex ===
	nodesMutex  sync.Mutex
	nodes       map[string]*t_node
	sortedNodes []*t_node
	// === nodesMutex ===

	// === datafeed ===
	nodeStatusFeedSource   *t_nodeStatusFeedSource
	nodeResourceFeedSource *t_nodeResourceFeedSource
	deviceStatusFeedSource *t_deviceStatusFeedSource
	streamStatusFeedSource *t_streamStatusFeedSource
	// === datafeed ===
}

func NewLiveCache(params *LiveCacheParams) (*LiveCache, error) {
	var err error

	lc := &LiveCache{}
	lc.context, lc.cancel = context.WithCancel(params.ParentContext)
	lc.logger, lc.logPrefix = params.ParentLogger, params.ParentLogPrefix + ".live"

	lc.evHub = params.EvHub

	lc.cron = cron.New(
		cron.WithLocation(params.Timezone),
		cron.WithSeconds(),
	)
	lc.cronJobs = make(map[string]cron.EntryID)

	lc.commBundle, err = lc.newCommBundle()
	if err != nil {
		return nil, err
	}

	lc.miscJobChan = make(chan *t_miscJob)

	lc.nodes       = make(map[string]*t_node)
	lc.sortedNodes = make([]*t_node, 0)

	return lc, nil
}

func (lc *LiveCache) Init() (error) {
	return nil
}

func (lc *LiveCache) Start() (error) {
	var err error

	lc.initNodeStatusFeedSource()
	lc.initNodeResourceFeedSource()
	lc.initDeviceStatusFeedSource()
	lc.initStreamStatusFeedSource()

	err = lc.commBundle.start()
	if err != nil { return err }

	lc.cron.Start()

	return nil
}

func (lc *LiveCache) PostStart() {
	go lc.startupRoutine()
}

func (lc *LiveCache) Stop() {
	lc.cron.Stop()

	_ = lc.commBundle.stop()
}

// func (lc *LiveCache) Stat() {}
