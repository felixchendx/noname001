package mdpservice

import (
	"encoding/json"

	baseTyping "noname001/app/base/typing"
)

type DeviceSnapshotRequest struct {
	ProviderNodeID string
	DeviceCode     string
}

type DeviceSnapshotReply struct {
	Status  string
	Message string

	Data    *baseTyping.BaseDeviceSnapshot
}

func (rcvr *DeviceReceiver) SendRequest_DeviceSnapshot(reqStruct *DeviceSnapshotRequest) (*DeviceSnapshotReply, error) {
	var err error
	
	deviceProviderID := DeviceProviderID(reqStruct.ProviderNodeID)

	reqBytes, err := json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_DeviceSnapshot - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD_DEVICE_SNAPSHOT + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(deviceProviderID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_DeviceSnapshot - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *DeviceSnapshotReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_DeviceSnapshot - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *DeviceProvider) recvRequest_DeviceSnapshot(reqJson string) (repJson string) {
	var reqStruct *DeviceSnapshotRequest

	err := json.Unmarshal([]byte(reqJson), &reqStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_DeviceSnapshot - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct DeviceSnapshotReply = DeviceSnapshotReply{}

	deviceSnapshot, messages := pvdr.dataHandler.ProvideDeviceSnapshot(reqStruct.DeviceCode)
	if messages.HasError() {
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil
	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = deviceSnapshot
	}

	repBytes, err := json.Marshal(repStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_DeviceSnapshot - req marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
