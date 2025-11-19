package mdpservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"noname001/logging"

	mdapi "noname001/dilemma/comm/zmqdep/mdp"
)

type StreamProviderParams struct {
	Context      context.Context
	Logger       *logging.WrappedLogger
	LogPrefix    string

	TempNodeID   string
	BrokerHost   string
	DataHandler  DataHandlerIntface

	Verbose      bool
	RetryBackoff time.Duration
}

type StreamProvider struct {
	context      context.Context
	cancel       context.CancelFunc
	logger       *logging.WrappedLogger
	logPrefix    string

	tempNodeID   string
	providerID   string
	brokerAddr   string
	dataHandler  DataHandlerIntface

	verbose      bool
	retryBackoff time.Duration

	mdworker     *mdapi.Mdwrk
}

func NewStreamProvider(params *StreamProviderParams) (*StreamProvider, error) {
	var err error
	
	pvdr := &StreamProvider{}
	pvdr.context, pvdr.cancel = context.WithCancel(params.Context)
	pvdr.logger = params.Logger
	pvdr.logPrefix = params.LogPrefix + ".strm.pvdr"

	pvdr.tempNodeID = params.TempNodeID
	pvdr.providerID = StreamProviderID(pvdr.tempNodeID)
	pvdr.brokerAddr = fmt.Sprintf("tcp://%s", params.BrokerHost)
	pvdr.dataHandler = params.DataHandler

	pvdr.verbose = params.Verbose
	pvdr.retryBackoff = params.RetryBackoff
	
	pvdr.mdworker = nil

	if err != nil {
		return nil, err
	}

	return pvdr, nil
}

func (pvdr *StreamProvider) AssignDataHandler(dataHandler DataHandlerIntface) {
	pvdr.dataHandler = dataHandler
}

func (pvdr *StreamProvider) Connect() (err error) {
	// TODO: idempotency

	if pvdr.dataHandler == nil {
		err = fmt.Errorf("Data handler is not provided.")
		pvdr.logger.Errorf("%s: %s", pvdr.logPrefix, err.Error())
		return err
	}

	go func() {
		// LimboLoop:
		for {
			var err error

			pvdr.mdworker, err = mdapi.NewMdwrk(pvdr.brokerAddr, pvdr.providerID, pvdr.verbose)
			if err != nil {
				// TODO: event channel
				pvdr.logger.Errorf("%s: new mdwrk err %s", pvdr.logPrefix, err.Error())
				time.Sleep(pvdr.retryBackoff)
				continue
			}

			var request, reply []string

			// RecvLoop:
			for {
				request, err = pvdr.mdworker.Recv(reply)
				if err != nil {
					// TODO: error for too long / unrecoverable err
					// TODO: ev channel
					pvdr.logger.Errorf("%s: recv err %s", pvdr.logPrefix, err.Error())
					time.Sleep(1 * time.Second)
					continue
				}

				reply = make([]string, 1, 1)

				if len(request) > 0 {
					reqParts := strings.Split(request[0], REQPARAM_DELIM)
					reqCode := reqParts[0]
					reqParam := ""
					if len(reqParts) >= 2 {
						reqParam = reqParts[1]
					}

					switch reqCode {
					case CMD_PING                : reply[0] = "pong"
					case CMD_SERVICE_INFO        : reply[0] = pvdr.recvRequest_ServiceInfo(reqParam)
					case CMD_STREAM_SNAPSHOT_LIST: reply[0] = pvdr.recvRequest_StreamSnapshotList(reqParam)
					case CMD_STREAM_SNAPSHOT     : reply[0] = pvdr.recvRequest_StreamSnapshot(reqParam)
					default:
						err = fmt.Errorf("unknown cmd '%s'", reqCode)
						pvdr.logger.Errorf("%s: ", pvdr.logPrefix, err.Error())
						reply[0] = SerializedErrorReply(err.Error())
					}
				}
			}
		}
	}()

	return
}

func (pvdr *StreamProvider) Disconnect() (err error) {
	return
}

// DO NOT USE, will cause panic
// func (sp *StreamProvider) Stop() {
// 	if sp.MDWorker != nil {
// 		sp.MDWorker.Close()
// 		sp.MDWorker = nil // cause panic further in

// 		sp.ccl()
// 	}
// }
