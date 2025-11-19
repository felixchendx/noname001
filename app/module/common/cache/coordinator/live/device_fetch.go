package live

import (
	"fmt"

	baseTyping "noname001/app/base/typing"

	deviceMdp "noname001/app/node/comm/device/mdpservice"
)

func (lc *LiveCache) fetchDeviceServiceInfo(nodeID string) (*deviceMdp.ServiceInfo, error) {
	req := &deviceMdp.ServiceInfoRequest{ProviderNodeID: nodeID}
	rep, err := lc.commBundle.deviceReceiver.SendRequest_ServiceInfo(req)
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

func (lc *LiveCache) fetchDeviceSnapshotList(nodeID string) ([]*baseTyping.BaseDeviceSnapshot, error) {
	req := &deviceMdp.DeviceSnapshotListRequest{ProviderNodeID: nodeID}
	rep, err := lc.commBundle.deviceReceiver.SendRequest_DeviceSnapshotList(req)
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

func (lc *LiveCache) fetchDeviceSnapshot(nodeID, deviceCode string) (*baseTyping.BaseDeviceSnapshot, error) {
	req := &deviceMdp.DeviceSnapshotRequest{ProviderNodeID: nodeID, DeviceCode: deviceCode}
	rep, err := lc.commBundle.deviceReceiver.SendRequest_DeviceSnapshot(req)
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
