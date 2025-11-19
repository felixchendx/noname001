package live

import (
	"sync"
	"time"

	streamTyping "noname001/app/base/typing/stream"
)

const (
	stream__external_activity_logs_limit = 20
	stream__internal_activity_logs_limit = 10

	stream__last_external_activity_threshold int64 = 60 * 3 // 3 mins

	stream__ttl int64 = 60 * 60 * 1 // 1 hour
)


func (lc *LiveCache) newStreamService() (*t_streamService) {
	return &t_streamService{
		streams: make(map[string]*t_stream),
		sortedStreams: make([]*t_stream, 0),

		opFlag  : false,
		opStatus: "",
	}
}

// -----------------------------------------------------------------------------
type t_streamService struct {
	// === streamsMutex ===
	streamsMutex  sync.Mutex
	streams       map[string]*t_stream
	sortedStreams []*t_stream
	// === streamsMutex ===

	opFlag   bool
	opStatus string
}

func (streamService *t_streamService) addNewStream(streamSnapshot *streamTyping.StreamSnapshot) (*t_stream) {
	stream := &t_stream{}
	stream.id = streamSnapshot.ID
	stream.code = streamSnapshot.Code

	stream.streamSnapshot = streamSnapshot
	stream.streamStatus = nil

	stream.lastInternalActivityAt = time.Now()
	stream.internalActivityLogs = make([]*t_internalActivityLog, 0, stream__internal_activity_logs_limit + 1)

	stream.lastExternalActivityAt = time.Now()
	stream.externalActivityLogs = make([]*t_externalActivityLog, 0, stream__external_activity_logs_limit + 1)

	stream.staleSince = zero_time
	stream.ttl = stream__ttl

	streamService.registerStream(stream)

	return stream
}

func (streamService *t_streamService) updateStreamData(stream *t_stream, streamSnapshot *streamTyping.StreamSnapshot) {
	stream.streamSnapshot = streamSnapshot

	stream.defunctAt = zero_time
	stream.defunctReason = ""

	stream.logExternalActivity("data_update", nil)
}

func (streamService *t_streamService) markStreamAsDefunct(stream *t_stream, reason string) {
	stream.defunctAt = time.Now()
	stream.defunctReason = reason

	stream.logExternalActivity("defunct", []string{reason})
}

func (streamService *t_streamService) checkStreamStaleness(stream *t_stream, checkTime time.Time) {
	if (checkTime.Unix() - stream.lastExternalActivityAt.Unix()) > stream__last_external_activity_threshold {
		if stream.staleSince.IsZero() {
			stream.staleSince = checkTime
		}
		stream.logInternalActivity("stale_check", "stale")

	} else {
		stream.staleSince = zero_time
		stream.logInternalActivity("stale_check", "ok")
	}
}

func (streamService *t_streamService) isStreamExpired(stream *t_stream, checkTime time.Time) (bool) {
	var expired bool = false

	if stream.staleSince.IsZero() {
		expired = false

	} else {
		expired = (checkTime.Unix() - stream.lastExternalActivityAt.Unix()) > stream.ttl
	}

	return expired
}

func (streamService *t_streamService) removeStream(stream *t_stream) {
	streamService.deregisterStream(stream)
}

// -----------------------------------------------------------------------------
type t_stream struct {
	id   string
	code string

	// === data holder ===
	streamSnapshot *streamTyping.StreamSnapshot
	streamStatus   *t_streamStatusInterpretation
	// === data holder ===

	// === generated ===
	sourceStream  string
	relayPathName string

	// === internal ===
	lastInternalActivityAt time.Time
	internalActivityLogs []*t_internalActivityLog

	lastExternalActivityAt time.Time
	externalActivityLogs []*t_externalActivityLog

	defunctAt     time.Time
	defunctReason string

	staleSince time.Time
	ttl int64
}
type t_streamStatusInterpretation struct {
	timestamp time.Time

	topLevelFail bool

	textualIndicator  string
	visualIndicator   string
	auditoryIndicator string
}

func (stream *t_stream) logInternalActivity(activity, result string) {
	var naw = time.Now()

	stream.lastInternalActivityAt = naw
	stream.internalActivityLogs = append(stream.internalActivityLogs, &t_internalActivityLog{
		ts: naw,
		activity: activity,
		result: result,
	})

	if len(stream.internalActivityLogs) > stream__internal_activity_logs_limit {
		stream.internalActivityLogs = stream.internalActivityLogs[1:len(stream.internalActivityLogs)]
	}
}

func (stream *t_stream) logExternalActivity(activity string, extra []string) {
	var naw = time.Now()

	stream.lastExternalActivityAt = naw
	stream.externalActivityLogs = append(stream.externalActivityLogs, &t_externalActivityLog{
		ts: naw,
		activity: activity,
		extra: extra,
	})

	if len(stream.externalActivityLogs) > stream__external_activity_logs_limit {
		stream.externalActivityLogs = stream.externalActivityLogs[1:len(stream.externalActivityLogs)]
	}
}
