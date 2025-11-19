package event

import (
	"context"
	"fmt"
	"time"

	"noname001/logging"
)

type BaseEventHubParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger
	LogPrefix     string
}
type BaseEventHub struct {
	Context   context.Context
	Cancel    context.CancelFunc
	Logger    *logging.WrappedLogger
	LogPrefix string

	sources map[string]*BaseEventSource
}

func NewBaseEventHub(params *BaseEventHubParams) (*BaseEventHub) {
	evHub := &BaseEventHub{}
	evHub.Context, evHub.Cancel = context.WithCancel(params.ParentContext)
	evHub.Logger, evHub.LogPrefix = params.Logger, params.LogPrefix + ".evHub"

	evHub.sources = make(map[string]*BaseEventSource)

	evHub._cleanupWorker()

	return evHub
}

func (evHub *BaseEventHub) NewEventSource(evSourceTree string) (error) {
	_, inMap := evHub.sources[evSourceTree]
	if inMap {
		err := fmt.Errorf("NewEventSource err, evSourceTree '%s' already used.", evSourceTree)
		evHub.Logger.Debugf("%s: %s", evHub.LogPrefix, err.Error())
		return err
	}

	evHub.sources[evSourceTree] = NewBaseEventSource(evHub.Context, evSourceTree)

	return nil
}

func (evHub *BaseEventHub) PublishToEventSource(evSourceTree string, ev any) {
	evSource, inMap := evHub.sources[evSourceTree]
	if inMap {
		evSource.Publish(ev)
	}
}

func (evHub *BaseEventHub) CloseEventSource(evSourceTree string) {
	evSource, inMap := evHub.sources[evSourceTree]
	if inMap {
		evSource.Close()
		delete(evHub.sources, evSourceTree)
	}
}

func (evHub *BaseEventHub) NewSubscription(evSourceTree string) (*BaseEventSubscription) {
	evSource, inMap := evHub.sources[evSourceTree]
	if !inMap { return nil }
	if evSource.closed { return nil }

	evSub := evSource.Subscribe()

	return evSub
}

func (evHub *BaseEventHub) CloseSubscription(evSub *BaseEventSubscription) {
	evSub.Unsubscribe()
}

func (evHub *BaseEventHub) _cleanupWorker() {
	go func() {
		var (
			cleanupTicker = time.NewTicker(60 * time.Second)
		)

		defer func() {
			cleanupTicker.Stop()
		}()

		workerLoop:
		for {
			select {
			case <- evHub.Context.Done():
				break workerLoop
	
			case <- cleanupTicker.C:
				for _, evSource := range evHub.sources {
					evSource._unsubCleanup()
				}
			}
		}
	}()
}
