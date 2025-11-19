package pubsub

const (
	// TODO: configurable base for different cluster
	SUBTREE__NODE_BASE     = "/sumthinkrando/node/"
	SUBTREE__NODE_EVENT    = "/sumthinkrando/node/event/"
	SUBTREE__NODE_LIVENESS = "/sumthinkrando/node/liveness/"
)

const (
	DEFAULT_TTL = "600" // 10 mins
)

const (
	HEADER_PAYLOAD_DELIM = "|><|"

	HEADER__EVENT    = "evnt"
	HEADER__LIVENESS = "lvns"
)
