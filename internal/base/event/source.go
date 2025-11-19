package event

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type BaseEventSource struct {
	context context.Context
	cancel  context.CancelFunc
	
	id string

	sourceChan chan any
	closed     bool     // TODO: there should be go channel status that indicate close or open... lookup doc later

	subscriptionsMutex sync.Mutex
	subscriptions      map[string]*BaseEventSubscription
}

func NewBaseEventSource(parentContext context.Context, evSourceTree string) (*BaseEventSource) {
	evSource := &BaseEventSource{}
	evSource.context, evSource.cancel = context.WithCancel(parentContext)

	evSource.id = evSourceTree
	evSource.sourceChan = make(chan any)
	evSource.closed = false
	evSource.subscriptions = make(map[string]*BaseEventSubscription)

	go func() {
		evSourceLoop:
		for {
			select {
			case <- evSource.context.Done():
				break evSourceLoop

			case ev, open := <- evSource.sourceChan:
				if !open { break evSourceLoop }

				// mutex on read ?
				for _, evSub := range evSource.subscriptions {
					if evSub.closed { continue }

					evSub.subChan <- ev
				}
			}
		}
	}()

	return evSource
}

func (evSource *BaseEventSource) Publish(ev any) {
	if evSource.closed { return }

	evSource.sourceChan <- ev
}

func (evSource *BaseEventSource) Subscribe() (*BaseEventSubscription) {
	evSub := &BaseEventSubscription{
		id: uuid.New().String(),

		subChan: make(chan any),

		unsubbed: false,
		closed  : evSource.closed,
	}

	evSource.subscriptionsMutex.Lock()
	evSource.subscriptions[evSub.id] = evSub
	evSource.subscriptionsMutex.Unlock()

	return evSub
}

func (evSource *BaseEventSource) Unsubscribe(evSub *BaseEventSubscription) {
	// mutex on read ?
	evSub, inMap := evSource.subscriptions[evSub.id]
	if inMap {
		evSub.Unsubscribe()
	}
}

func (evSource *BaseEventSource) Close() {
	for _, evSub := range evSource.subscriptions {
		if evSub.closed { continue }

		evSub.closed = true
		// close(evSub.subChan)
	}

	evSource.subscriptionsMutex.Lock()
	evSource.subscriptions = make(map[string]*BaseEventSubscription)
	evSource.subscriptionsMutex.Unlock()

	evSource.closed = true
	// close(evSource.sourceChan)
	evSource.cancel()
}

func (evSource *BaseEventSource) _unsubCleanup() {
	cleanupList := make([]string, 0)

	for _, evSub := range evSource.subscriptions {
		if evSub.closed {
			cleanupList = append(cleanupList, evSub.id)
		}
	} 

	if len(cleanupList) > 0 {
		evSource.subscriptionsMutex.Lock()
		for _, evSubID := range cleanupList {
			delete(evSource.subscriptions, evSubID)
		}
		evSource.subscriptionsMutex.Unlock()
	}
}
