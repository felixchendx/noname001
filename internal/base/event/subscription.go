package event

type BaseEventSubscription struct {
	id string

	subChan chan any

	unsubbed bool
	closed   bool
}

func (evSub *BaseEventSubscription) Channel() (chan any) {
	return evSub.subChan
}

func (evSub *BaseEventSubscription) Unsubscribe() {
	evSub.unsubbed = true
	evSub.closed   = true
	// close(evSub.subChan)
}
