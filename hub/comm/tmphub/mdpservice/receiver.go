package mdpservice

import (
	"fmt"

	"noname001/logging"

	mdapi "noname001/dilemma/comm/zmqdep/mdp"
)

type TmpHubReceiverParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BrokerHost string

	Verbose    bool
}
type TmpHubReceiver struct {
	logger     *logging.WrappedLogger
	logPrefix  string

	brokerAddr string

	mdclient   *mdapi.Mdcli
}

func NewTmpHubReceiver(params *TmpHubReceiverParams) (*TmpHubReceiver, error) {
	var err error
	
	rcvr := &TmpHubReceiver{}
	rcvr.logger = params.Logger
	rcvr.logPrefix = params.LogPrefix + ".tmphub.rcvr"

	rcvr.brokerAddr = fmt.Sprintf("tcp://%s", params.BrokerHost)

	rcvr.mdclient, err = mdapi.NewMdcli(rcvr.brokerAddr, params.Verbose)
	if err != nil {
		rcvr.logger.Errorf("%s: new receiver err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	return rcvr, nil
}

func (rcvr *TmpHubReceiver) Connect() (err error) {
	return
}

func (rcvr *TmpHubReceiver) Disconnect() (err error) {
	return
}
