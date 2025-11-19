package event

var (
	evHubInstance *EventHub
)

func EventHubInstance() (*EventHub) {
	return evHubInstance
}
