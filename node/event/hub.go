package event

import (
	"context"

	"noname001/logging"
)

type EventHubParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger
	LogPrefix     string
}
type EventHub struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	nodeEventSource *t_nodeEventSource
}

func NewEventHub(params *EventHubParams) (*EventHub) {
	evHub := &EventHub{}
	evHub.context, evHub.cancel = context.WithCancel(params.ParentContext)
	evHub.logger, evHub.logPrefix = params.Logger, params.LogPrefix + ".evHub"

	return evHub
}

func (evHub *EventHub) Open() {
	evHub.nodeEventSource = evHub.newNodeEventSource()
}

func (evHub *EventHub) Close() {
}
