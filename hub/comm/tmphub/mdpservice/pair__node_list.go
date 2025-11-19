package mdpservice

import (
	"encoding/json"
)

type NodeListRequest struct {
	// TODO: service filter
}

type NodeListReply struct {
	Status  string
	Message string

	Data    []*NodeInfo
}
type NodeInfo struct {
	NodeID      string
	ServiceList []string
}

func (rcvr *TmpHubReceiver) SendRequest_NodeList(reqStruct *NodeListRequest) (*NodeListReply, error) {
	var err error
	
	reqBytes, err := json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_NodeList - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD_NODE_LIST + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(PROVIDER_ID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_NodeList - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *NodeListReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_NodeList - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *TmpHubProvider) recvRequest_NodeList(reqJson string) (repJson string) {
	var reqStruct *NodeListRequest

	err := json.Unmarshal([]byte(reqJson), &reqStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_NodeList - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct NodeListReply = NodeListReply{}

	nodeList, messages := pvdr.dataHandler.ProvideNodeList()
	if messages.HasError() {
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil
	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = nodeList
	}

	repBytes, err := json.Marshal(repStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_NodeList - req marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
