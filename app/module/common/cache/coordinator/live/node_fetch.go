package live

import (
	"fmt"

	nodeTyping "noname001/node/base/typing"

	tmphubMdp "noname001/hub/comm/tmphub/mdpservice"
	nodeMdp   "noname001/node/comm/mdpservice"
)

func (lc *LiveCache) fetchNodeInfoList() ([]*tmphubMdp.NodeInfo, error) {
	rep, err := lc.commBundle.tmpHubReceiver.SendRequest_NodeList(&tmphubMdp.NodeListRequest{})
	if err != nil {
		// TODO: evChan
		lc.logger.Errorf("%s: %s", lc.logPrefix, err.Error())
		return nil, err
	}
	if rep.Status == "error" {
		// TODO: evChan
		lc.logger.Errorf("%s: %s", lc.logPrefix, rep.Message)
		return nil, fmt.Errorf(rep.Message)
	}

	return rep.Data, nil
}

func (lc *LiveCache) fetchNodeSnapshot(nodeID string) (*nodeTyping.BaseNodeSnapshot, error) {
	req := &nodeMdp.NodeSnapshotRequest{NodeID: nodeID}
	rep, err := lc.commBundle.nodeReceiver.SendRequest_NodeSnapshot(req)
	if err != nil {
		// TODO: evChan
		lc.logger.Errorf("%s: node '%s' - %s", lc.logPrefix, nodeID, err.Error())
		return nil, err
	}
	if rep.Status == "error" {
		// TODO: evChan
		lc.logger.Errorf("%s: node '%s' - %s", lc.logPrefix, nodeID, rep.Message)
		return nil, fmt.Errorf(rep.Message)
	}

	return rep.Data, nil
}

func (lc *LiveCache) fetchNodeResource(nodeID string) (*nodeTyping.TempNodeSystemResourceSummary, error) {
	req := &nodeMdp.NodeResourceRequest{NodeID: nodeID}
	rep, err := lc.commBundle.nodeReceiver.SendRequest_NodeResource(req)
	if err != nil {
		// TODO: evChan
		lc.logger.Errorf("%s: node '%s' - %s", lc.logPrefix, nodeID, err.Error())
		return nil, err
	}
	if rep.Status == "error" {
		// TODO: evChan
		lc.logger.Errorf("%s: node '%s' - %s", lc.logPrefix, nodeID, rep.Message)
		return nil, fmt.Errorf(rep.Message)
	}

	return rep.Data, nil
}
