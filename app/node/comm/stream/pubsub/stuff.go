package pubsub

const (
	// TODO: configurable base for different cluster
	SUBTREE__STREAM_BASE     = "/sumthinkrando/stream/"
	SUBTREE__STREAM_EVENT    = "/sumthinkrando/stream/event/"
	SUBTREE__STREAM_LIVENESS = "/sumthinkrando/stream/liveness/"
)

const (
	DEFAULT_TTL = "600" // 10 mins
)

const (
	HEADER_PAYLOAD_DELIM = "|><|"

	HEADER__EVENT    = "evnt"
	HEADER__LIVENESS = "lvns"
)
