package mdpservice

import (
	"fmt"

	"noname001/logging"

	mdapi "noname001/dilemma/comm/zmqdep/mdp"
)

type NodeReceiverParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BrokerHost string

	Verbose    bool
}
type NodeReceiver struct {
	logger     *logging.WrappedLogger
	logPrefix  string

	brokerAddr string

	mdclient   *mdapi.Mdcli
}

func NewNodeReceiver(params *NodeReceiverParams) (*NodeReceiver, error) {
	var err error
	
	rcvr := &NodeReceiver{}
	rcvr.logger = params.Logger
	rcvr.logPrefix = params.LogPrefix + ".node.rcvr"

	rcvr.brokerAddr = fmt.Sprintf("tcp://%s", params.BrokerHost)

	rcvr.mdclient, err = mdapi.NewMdcli(rcvr.brokerAddr, params.Verbose)
	if err != nil {
		rcvr.logger.Errorf("%s: new receiver err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	return rcvr, nil
}

func (rcvr *NodeReceiver) Connect() (err error) {
	return
}

func (rcvr *NodeReceiver) Disconnect() (err error) {
	return
}
