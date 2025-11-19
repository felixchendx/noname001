package live

import (
	"context"

	"noname001/logging"

	streamTyping "noname001/app/base/typing/stream"

	streamEv    "noname001/app/module/feature/stream/event"
	streamStore "noname001/app/module/feature/stream/store"
	streamComm  "noname001/app/module/feature/stream/comm"
)

type LiveStreamParams struct {
	ParentContext   context.Context
	ParentLogger    *logging.WrappedLogger
	ParentLogPrefix string

	EvHub      *streamEv.EventHub
	Store      *streamStore.Store
	CommBundle *streamComm.CommBundle

	StreamID string
}

type LiveStream struct {
	context    context.Context
	cancel     context.CancelFunc
	logger     *logging.WrappedLogger
	logPrefix  string

	evHub      *streamEv.EventHub
	store      *streamStore.Store
	commBundle *streamComm.CommBundle

	id   string
	code string

	isDestroying     bool
	execChan         chan t_execCode
	streamerInstance *ffmpegStreamer
	logHistories     []*t_internalLog

	inputs  []inputIntface
	outputs []*outputDestinationInternalMediaServer

	pDat streamTyping.StreamPersistenceData
	lDat streamTyping.StreamLiveData
}

func NewLiveStream(params *LiveStreamParams) (*LiveStream) {
	liveStream := &LiveStream{}
	liveStream.context, liveStream.cancel = context.WithCancel(params.ParentContext)
	liveStream.logger, liveStream.logPrefix = params.ParentLogger, params.ParentLogPrefix + ".live"
	liveStream.evHub = params.EvHub
	liveStream.store = params.Store
	liveStream.commBundle = params.CommBundle

	liveStream.id   = params.StreamID
	liveStream.code = ""

	liveStream.isDestroying     = false
	liveStream.execChan         = make(chan t_execCode, 16)
	liveStream.streamerInstance = nil
	liveStream.logHistories     = make([]*t_internalLog, 0)

	liveStream.lDat.State     = streamTyping.LIVE_STATE__NEW
	liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE

	go liveStream.executor()
	go liveStream.watcher()
	go liveStream.evListener()

	return liveStream
}

func (liveStream *LiveStream) Init() {
	if liveStream.isDestroying { return }

	liveStream.execChan <- exec_init
}

func (liveStream *LiveStream) Reload() {
	if liveStream.isDestroying { return }

	liveStream.execChan <- exec_reload
}

func (liveStream *LiveStream) Destroy() {
	if liveStream.isDestroying { return }

	liveStream.isDestroying = true

	liveStream.execChan <- exec_destroy
}

func (liveStream *LiveStream) Snapshot() (*streamTyping.StreamSnapshot) {
	return &streamTyping.StreamSnapshot{
		ID  : liveStream.id,
		Code: liveStream.code,

		Persistence: liveStream.pDat,
		Live       : liveStream.lDat,
	}
}

func (liveStream *LiveStream) Code() (streamCode string) {
	return liveStream.code
}

func (liveStream *LiveStream) PersistenceData() (streamTyping.StreamPersistenceData) {
	return liveStream.pDat
}

// temp, group is to be dismissed
func (liveStream *LiveStream) DestroyIfBelongToThisGroup(groupID string) (bool) {
	if liveStream.isDestroying { return true }

	if liveStream.pDat.GroupID != groupID { return false }

	liveStream.isDestroying = true

	liveStream.execChan <- exec_destroy

	return true
}
