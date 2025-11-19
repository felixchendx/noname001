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

	cachedNodeEventSource   *t_cachedNodeEventSource
	cachedDeviceEventSource *t_cachedDeviceEventSource
	cachedStreamEventSource *t_cachedStreamEventSource
}

func NewEventHub(params *EventHubParams) (*EventHub) {
	evHub := &EventHub{}
	evHub.context, evHub.cancel = context.WithCancel(params.ParentContext)
	evHub.logger, evHub.logPrefix = params.Logger, params.LogPrefix + ".evHub"

	return evHub
}

func (evHub *EventHub) Open() {
	evHub.cachedNodeEventSource   = evHub.newCachedNodeEventSource()
	evHub.cachedDeviceEventSource = evHub.newCachedDeviceEventSource()
	evHub.cachedStreamEventSource = evHub.newCachedStreamEventSource()
}

func (evHub *EventHub) Close() {
}
