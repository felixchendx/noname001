package hub

import (
	"strings"

	"noname001/dilemma/comm"

	"noname001/app/base/messaging"

	tmphubMdp "noname001/hub/comm/tmphub/mdpservice"
)

// temp ?
type tmphubProvider_dataHandler_t struct {
	commBroker *comm.BrokerWrapper
}

// ============ VVV conform to tmphubMdp.DataHandlerIntface VVV ============= //
func (dataHandler *tmphubProvider_dataHandler_t) ProvideNodeList() ([]*tmphubMdp.NodeInfo, *messaging.Messages) {
	messages := messaging.NewMessages()
	nodeList := make([]*tmphubMdp.NodeInfo, 0)

	// TODO: precompute to cache
	serviceList := dataHandler.commBroker.RetrieveServiceList()
	for k, _ := range serviceList {
		parts := strings.Split(k, "::")
		if len(parts) != 2 {
			continue
		}

		nodeID, serviceCode := parts[0], parts[1]

		if nodeID == "TMPHUB" { continue }

		var currNodeInfo *tmphubMdp.NodeInfo
		for _, item := range nodeList {
			if nodeID == item.NodeID {
				currNodeInfo = item
				break
			}
		}

		if currNodeInfo == nil {
			currNodeInfo = &tmphubMdp.NodeInfo{
				NodeID: nodeID,
				ServiceList: make([]string, 0),
			}
			nodeList = append(nodeList, currNodeInfo)
		}

		currNodeInfo.ServiceList = append(currNodeInfo.ServiceList, serviceCode)
	}

	return nodeList, messages 
}
// ============ ^^^ conform to tmphubMdp.DataHandlerIntface ^^^ ============= //
