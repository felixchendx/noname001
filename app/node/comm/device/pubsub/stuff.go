package pubsub

const (
	// TODO: configurable base for different cluster
	SUBTREE__DEVICE_BASE     = "/sumthinkrando/device/"
	SUBTREE__DEVICE_EVENT    = "/sumthinkrando/device/event/"
)

const (
	DEFAULT_TTL = "600" // 10 mins
)

const (
	HEADER_PAYLOAD_DELIM = "|><|"

	HEADER__EVENT    = "evnt"
)
