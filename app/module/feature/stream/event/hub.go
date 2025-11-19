package event

import (
	"context"

	"noname001/logging"
)

type EventHubParams struct {
	ParentContext   context.Context
	Logger          *logging.WrappedLogger
	LogPrefix       string
}
type EventHub struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	streamProfileEventSource *t_streamProfileEventSource
	streamGroupEventSource   *t_streamGroupEventSource
	streamItemEventSource    *t_streamItemEventSource

	liveStreamEventSource    *t_liveStreamEventSource
}

func NewEventHub(params *EventHubParams) (*EventHub) {
	evHub := &EventHub{}
	evHub.context, evHub.cancel = context.WithCancel(params.ParentContext)
	evHub.logger, evHub.logPrefix = params.Logger, params.LogPrefix + ".evHub"

	return evHub
}

func (evHub *EventHub) Open() {
	evHub.streamProfileEventSource = evHub.newStreamProfileEventSource()
	evHub.streamGroupEventSource   = evHub.newStreamGroupEventSource()
	evHub.streamItemEventSource    = evHub.newStreamItemEventSource()

	evHub.liveStreamEventSource    = evHub.newLiveStreamEventSource()
}

func (evHub *EventHub) Close() {
}
