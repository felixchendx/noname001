package commconf

// cross node comm config
// rehash the concept, so that folder organization reflects the intended concept
// and after that, reorg comm stuffs

var (
	TempID   string
	TempName string
	TempIPs  []string // TODO: this needs to be dynamic, for re-broadcast
	TempBrokerServerHost string
	TempSnapshotServerHost string
	TempPublisherServerHost string
	TempCollectorServerHost string
	TempCommVerbose bool
)

func ID() (string) {
	return TempID
}

func Name() (string) {
	return TempName
}

func IPs() ([]string) {
	return TempIPs
}

func BrokerServerHost() (string) {
	return TempBrokerServerHost
}

func SnapshotServerHost() (string) {
	return TempSnapshotServerHost
}

func PublisherServerHost() (string) {
	return TempPublisherServerHost
}

func CollectorServerHost() (string) {
	return TempCollectorServerHost
}

func CommVerbose() (bool) {
	return TempCommVerbose
}
