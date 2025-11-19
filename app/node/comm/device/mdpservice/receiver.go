package mdpservice

import (
	"fmt"

	"noname001/logging"

	mdapi "noname001/dilemma/comm/zmqdep/mdp"
)

type DeviceReceiverParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BrokerHost string

	Verbose    bool
}
type DeviceReceiver struct {
	logger     *logging.WrappedLogger
	logPrefix  string

	brokerAddr string

	mdclient   *mdapi.Mdcli
}

func NewDeviceReceiver(params *DeviceReceiverParams) (*DeviceReceiver, error) {
	var err error
	
	rcvr := &DeviceReceiver{}
	rcvr.logger = params.Logger
	rcvr.logPrefix = params.LogPrefix + ".dvc.rcvr"

	rcvr.brokerAddr = fmt.Sprintf("tcp://%s", params.BrokerHost)

	rcvr.mdclient, err = mdapi.NewMdcli(rcvr.brokerAddr, params.Verbose)
	if err != nil {
		rcvr.logger.Errorf("%s: new receiver err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	return rcvr, nil
}

func (rcvr *DeviceReceiver) Connect() (err error) {
	return
}

func (rcvr *DeviceReceiver) Disconnect() (err error) {
	return
}

func (rcvr *DeviceReceiver) GetServiceCode() (string) {
	return SERVICE_CODE
}
