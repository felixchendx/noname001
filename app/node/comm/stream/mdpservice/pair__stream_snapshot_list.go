package mdpservice

import (
	"encoding/json"

	streamTyping "noname001/app/base/typing/stream"
)

type StreamSnapshotListRequest struct {
	ProviderNodeID string
}

type StreamSnapshotListReply struct {
	Status  string
	Message string

	Data    []*streamTyping.StreamSnapshot
}

func (rcvr *StreamReceiver) SendRequest_StreamSnapshotList(reqStruct *StreamSnapshotListRequest) (*StreamSnapshotListReply, error) {
	var err error

	streamProviderID := StreamProviderID(reqStruct.ProviderNodeID)

	reqBytes, err := json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_StreamSnapshotList - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD_STREAM_SNAPSHOT_LIST + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(streamProviderID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_StreamSnapshotList - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *StreamSnapshotListReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_StreamSnapshotList - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *StreamProvider) recvRequest_StreamSnapshotList(reqJson string) (repJson string) {
	var reqStruct *StreamSnapshotListRequest

	err := json.Unmarshal([]byte(reqJson), &reqStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_StreamSnapshotList - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct StreamSnapshotListReply = StreamSnapshotListReply{}

	streamSnapshotList, messages := pvdr.dataHandler.ProvideStreamSnapshotList()
	if messages.HasError() {
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil
	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = streamSnapshotList
	}

	repBytes, err := json.Marshal(repStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_StreamSnapshotList - req marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
