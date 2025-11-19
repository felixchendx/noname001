package mdpservice

import (
	"encoding/json"

	nodeTyping "noname001/node/base/typing"
)

// fukken temp

type NodeResourceRequest struct {
	NodeID string
}

type NodeResourceReply struct {
	Status  string
	Message string

	Data    *nodeTyping.TempNodeSystemResourceSummary
}

func (rcvr *NodeReceiver) SendRequest_NodeResource(reqStruct *NodeResourceRequest) (*NodeResourceReply, error) {
	var err error

	nodeProviderID := nodeProviderID(reqStruct.NodeID)

	var reqBytes []byte
	reqBytes, err = json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_NodeResource - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD__NODE_RESOURCE + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(nodeProviderID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_NodeResource - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *NodeResourceReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_NodeResource - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *NodeProvider) recvRequest_NodeResource(reqJSON string) (repJSON string) {
	var err error
	var reqStruct *NodeResourceRequest

	err = json.Unmarshal([]byte(reqJSON), &reqStruct)
	if err != nil {
		// TODO
		pvdr.logger.Errorf("%s: recvRequest_NodeResource - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct NodeResourceReply = NodeResourceReply{}
	nodeResource, messages := pvdr.dataHandler.ProvideNodeResource()
	if messages.HasError() {
		// TODO
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil

	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = nodeResource
	}

	var repBytes []byte
	repBytes, err = json.Marshal(repStruct)
	if err != nil {
		// TODO
		pvdr.logger.Errorf("%s: recvRequest_NodeResource - rep marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
