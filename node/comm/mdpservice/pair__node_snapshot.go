package mdpservice

import (
	"encoding/json"

	nodeTyping "noname001/node/base/typing"
)

type NodeSnapshotRequest struct {
	NodeID string
}

type NodeSnapshotReply struct {
	Status  string
	Message string

	Data    *nodeTyping.BaseNodeSnapshot
}

func (rcvr *NodeReceiver) SendRequest_NodeSnapshot(reqStruct *NodeSnapshotRequest) (*NodeSnapshotReply, error) {
	var err error

	nodeProviderID := nodeProviderID(reqStruct.NodeID)

	var reqBytes []byte
	reqBytes, err = json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_NodeSnapshot - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD__NODE_SNAPSHOT + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(nodeProviderID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_NodeSnapshot - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *NodeSnapshotReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_NodeSnapshot - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *NodeProvider) recvRequest_NodeSnapshot(reqJSON string) (repJSON string) {
	var err error
	var reqStruct *NodeSnapshotRequest

	err = json.Unmarshal([]byte(reqJSON), &reqStruct)
	if err != nil {
		// TODO
		pvdr.logger.Errorf("%s: recvRequest_NodeSnapshot - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct NodeSnapshotReply = NodeSnapshotReply{}
	nodeSnapshot, messages := pvdr.dataHandler.ProvideNodeSnapshot()
	if messages.HasError() {
		// TODO
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil

	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = nodeSnapshot
	}

	var repBytes []byte
	repBytes, err = json.Marshal(repStruct)
	if err != nil {
		// TODO
		pvdr.logger.Errorf("%s: recvRequest_NodeSnapshot - rep marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
