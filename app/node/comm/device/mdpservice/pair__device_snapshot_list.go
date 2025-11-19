package mdpservice

import (
	"encoding/json"

	baseTyping "noname001/app/base/typing"
)

type DeviceSnapshotListRequest struct {
	ProviderNodeID string
	// TODO StateFilter
}

type DeviceSnapshotListReply struct {
	Status  string
	Message string

	Data    []*baseTyping.BaseDeviceSnapshot
}

func (rcvr *DeviceReceiver) SendRequest_DeviceSnapshotList(reqStruct *DeviceSnapshotListRequest) (*DeviceSnapshotListReply, error) {
	var err error

	deviceProviderID := DeviceProviderID(reqStruct.ProviderNodeID)

	reqBytes, err := json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_DeviceSnapshotList - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD_DEVICE_SNAPSHOT_LIST + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(deviceProviderID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_DeviceSnapshotList - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *DeviceSnapshotListReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_DeviceSnapshotList - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *DeviceProvider) recvRequest_DeviceSnapshotList(reqJson string) (repJson string) {
	var reqStruct *DeviceSnapshotListRequest

	err := json.Unmarshal([]byte(reqJson), &reqStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_DeviceSnapshotList - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct DeviceSnapshotListReply = DeviceSnapshotListReply{}

	deviceSnapshotList, messages := pvdr.dataHandler.ProvideDeviceSnapshotList()
	if messages.HasError() {
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil
	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = deviceSnapshotList
	}

	repBytes, err := json.Marshal(repStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_DeviceSnapshotList - req marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
