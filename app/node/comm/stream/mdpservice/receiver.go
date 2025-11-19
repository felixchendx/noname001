package mdpservice

import (
	"fmt"

	"noname001/logging"

	mdapi "noname001/dilemma/comm/zmqdep/mdp"
)

type StreamReceiverParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BrokerHost string

	Verbose    bool
}
type StreamReceiver struct {
	logger     *logging.WrappedLogger
	logPrefix  string

	brokerAddr string

	mdclient   *mdapi.Mdcli
}

func NewStreamReceiver(params *StreamReceiverParams) (*StreamReceiver, error) {
	var err error
	
	rcvr := &StreamReceiver{}
	rcvr.logger = params.Logger
	rcvr.logPrefix = params.LogPrefix + ".strm.rcvr"

	rcvr.brokerAddr = fmt.Sprintf("tcp://%s", params.BrokerHost)

	rcvr.mdclient, err = mdapi.NewMdcli(rcvr.brokerAddr, params.Verbose)
	if err != nil {
		rcvr.logger.Errorf("%s: new receiver err %s", rcvr.logPrefix, err.Error())
		return nil, err
	}

	return rcvr, nil
}

func (rcvr *StreamReceiver) Connect() (err error) {
	return
}

func (rcvr *StreamReceiver) Disconnect() (err error) {
	return
}

func (rcvr *StreamReceiver) GetServiceCode() (string) {
	return SERVICE_CODE
}
