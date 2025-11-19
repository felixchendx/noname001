package event

import (
	"context"

	"noname001/logging"

	baseEv "noname001/internal/base/event"
)

type EventHubParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger
	LogPrefix     string
}
type EventHub struct {
	*baseEv.BaseEventHub
}
type EventSource       = baseEv.BaseEventSource
type EventSubscription = baseEv.BaseEventSubscription

func NewEventHub(params *EventHubParams) (*EventHub) {
	baseEvHub := baseEv.NewBaseEventHub(&baseEv.BaseEventHubParams{
		ParentContext: params.ParentContext,
		Logger: params.Logger, LogPrefix: params.LogPrefix,
	})

	evHub := &EventHub{
		baseEvHub,
	}
	evHubInstance = evHub

	evHub.Logger.Infof("[%s] initialized", evHub.LogPrefix)
	return evHub
}

func (evHub *EventHub) Open() (error) {
	var err error

	err = evHub.NewEventSource(EVTREE__BACKGROUND)
	if err != nil {
		return err
	}

	evHub.bgEventListener()

	evHub.Logger.Infof("[%s] started", evHub.LogPrefix)
	return nil
}

func (evHub *EventHub) Close() {
	evHub.CloseEventSource(EVTREE__BACKGROUND)

	evHub.Cancel()

	evHub.Logger.Infof("[%s] stopped", evHub.LogPrefix)
}
