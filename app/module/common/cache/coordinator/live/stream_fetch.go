package live

import (
	"fmt"

	streamTyping "noname001/app/base/typing/stream"

	streamMdp "noname001/app/node/comm/stream/mdpservice"
)

func (lc *LiveCache) fetchStreamServiceInfo(nodeID string) (*streamMdp.ServiceInfo, error) {
	req := &streamMdp.ServiceInfoRequest{ProviderNodeID: nodeID}
	rep, err := lc.commBundle.streamReceiver.SendRequest_ServiceInfo(req)
	if err != nil {
		// TODO:
		lc.logger.Errorf("%s: %s", lc.logPrefix, err.Error())
		return nil, err
	}
	if rep.Status == "error" {
		// TODO:
		lc.logger.Errorf("%s: %s", lc.logPrefix, rep.Message)
		return nil, fmt.Errorf(rep.Message)
	}

	return rep.Data, nil
}

func (lc *LiveCache) fetchStreamSnapshotList(nodeID string) ([]*streamTyping.StreamSnapshot, error) {
	req := &streamMdp.StreamSnapshotListRequest{ProviderNodeID: nodeID}
	rep, err := lc.commBundle.streamReceiver.SendRequest_StreamSnapshotList(req)
	if err != nil {
		// TODO:
		lc.logger.Errorf("%s: %s", lc.logPrefix, err.Error())
		return nil, err
	}
	if rep.Status == "error" {
		// TODO:
		lc.logger.Errorf("%s: %s", lc.logPrefix, rep.Message)
		return nil, fmt.Errorf(rep.Message)
	}

	return rep.Data, nil
}

func (lc *LiveCache) fetchStreamSnapshot(nodeID, streamCode string) (*streamTyping.StreamSnapshot, error) {
	req := &streamMdp.StreamSnapshotRequest{ProviderNodeID: nodeID, StreamCode: streamCode}
	rep, err := lc.commBundle.streamReceiver.SendRequest_StreamSnapshot(req)
	if err != nil {
		// TODO:
		lc.logger.Errorf("%s: %s", lc.logPrefix, err.Error())
		return nil, err
	}
	if rep.Status == "error" {
		// TODO:
		lc.logger.Errorf("%s: %s", lc.logPrefix, rep.Message)
		return nil, fmt.Errorf(rep.Message)
	}

	return rep.Data, nil
}
