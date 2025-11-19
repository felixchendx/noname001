package intface

type StreamServiceProviderIntface interface {
	// hmm... cumbersome
}

type StreamEventProviderIntface interface {
	SubscribeToLiveStreamEvent() (*LiveStreamEventSubscription)
	UnsubscribeFromLiveStreamEvent(*LiveStreamEventSubscription)
}

func EventProvider() (StreamEventProviderIntface) {
	return streamEventProvider
}

// === ^^^ for those that uses     ^^^ ===
// =======================================
// === VVV for those that provides VVV ===

var (
	streamServiceProvider StreamServiceProviderIntface
	streamEventProvider   StreamEventProviderIntface
)

func AssignStreamEventProvider(_something StreamEventProviderIntface) {
	streamEventProvider = _something
}
