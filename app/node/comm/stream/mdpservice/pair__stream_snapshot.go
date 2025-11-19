package mdpservice

import (
	"encoding/json"

	streamTyping "noname001/app/base/typing/stream"
)

type StreamSnapshotRequest struct {
	ProviderNodeID string
	StreamCode     string
}

type StreamSnapshotReply struct {
	Status  string
	Message string

	Data    *streamTyping.StreamSnapshot
}

func (rcvr *StreamReceiver) SendRequest_StreamSnapshot(reqStruct *StreamSnapshotRequest) (*StreamSnapshotReply, error) {
	var err error
	
	streamProviderID := StreamProviderID(reqStruct.ProviderNodeID)

	reqBytes, err := json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_StreamSnapshot - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD_STREAM_SNAPSHOT + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(streamProviderID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_StreamSnapshot - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *StreamSnapshotReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_StreamSnapshot - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *StreamProvider) recvRequest_StreamSnapshot(reqJson string) (repJson string) {
	var reqStruct *StreamSnapshotRequest

	err := json.Unmarshal([]byte(reqJson), &reqStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_StreamSnapshot - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct StreamSnapshotReply = StreamSnapshotReply{}

	streamSnapshot, messages := pvdr.dataHandler.ProvideStreamSnapshot(reqStruct.StreamCode)
	if messages.HasError() {
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil
	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = streamSnapshot
	}

	repBytes, err := json.Marshal(repStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_StreamSnapshot - req marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
