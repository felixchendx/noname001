package event

// temporary home to sidestep circular import
// rehash the concept, so that folder organization reflects the intended concept
// need to add capabilities for event filtering + forwarding
// so that it's easier to isolate event hub instance depending on needs
// and after that, reorg

var (
	evHubInstance *EventHub
)

func LocalEventHub() (*EventHub) {
	return evHubInstance
}
