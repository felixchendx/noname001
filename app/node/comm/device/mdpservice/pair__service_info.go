package mdpservice

import (
	"encoding/json"
)

type ServiceInfoRequest struct {
	ProviderNodeID string
}

type ServiceInfoReply struct {
	Status  string
	Message string

	Data    *ServiceInfo
}

type ServiceInfo struct {
	Placeholder string
}

func (rcvr *DeviceReceiver) SendRequest_ServiceInfo(reqStruct *ServiceInfoRequest) (*ServiceInfoReply, error) {
	var err error
	
	deviceProviderID := DeviceProviderID(reqStruct.ProviderNodeID)

	reqBytes, err := json.Marshal(reqStruct)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_ServiceInfo - req marshal err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	reqPayload := CMD_SERVICE_INFO + REQPARAM_DELIM + string(reqBytes)

	var reply []string
	reply, err = rcvr.mdclient.Send(deviceProviderID, reqPayload)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: SendRequest_ServiceInfo - send err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	var repStruct *ServiceInfoReply
	if len(reply) > 0 {
		err = json.Unmarshal([]byte(reply[0]), &repStruct)
		if err != nil {
			// TODO
			rcvr.logger.Errorf("%s: SendRequest_ServiceInfo - rep marshal err %s", rcvr.logPrefix, err.Error())
			return nil, err
		}
	}

	return repStruct, nil
}

func (pvdr *DeviceProvider) recvRequest_ServiceInfo(reqJson string) (repJson string) {
	var reqStruct *ServiceInfoRequest

	err := json.Unmarshal([]byte(reqJson), &reqStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_ServiceInfo - req unmarshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	var repStruct ServiceInfoReply = ServiceInfoReply{}

	serviceInfo, messages := pvdr.dataHandler.ProvideServiceInfo()
	if messages.HasError() {
		repStruct.Status = "error"
		repStruct.Message = messages.FirstErrorMessageString()
		repStruct.Data = nil
	} else {
		repStruct.Status = "ok"
		repStruct.Message = ""
		repStruct.Data = serviceInfo
	}

	repBytes, err := json.Marshal(repStruct)
	if err != nil {
		pvdr.logger.Errorf("%s: recvRequest_ServiceInfo - req marshal err %s", pvdr.logPrefix, err.Error())
		return SerializedErrorReply(err.Error())
	}

	return string(repBytes)
}
