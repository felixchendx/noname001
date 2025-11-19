package node

import (
	"noname001/config/rawconfig"

	"noname001/node/commconf"
)

// TODO: BIG FAT TEMP
func (node *Node_t) injectCommConf(cfg *rawconfig.ConfigRoot) {
	commconf.TempID = node.id
	commconf.TempName = node.name

	commconf.TempBrokerServerHost = cfg.Node.BrokerServerHost
	commconf.TempSnapshotServerHost = cfg.Node.SnapshotServerHost
	commconf.TempPublisherServerHost = cfg.Node.PublisherServerHost
	commconf.TempCollectorServerHost = cfg.Node.CollectorServerHost
	commconf.TempCommVerbose = cfg.Node.CommVerbose
}
