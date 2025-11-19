package hub

import (
	"context"
	"fmt"

	"noname001/dilemma/comm"

	"noname001/logging"
	"noname001/config/rawconfig"

	tmphubMdp "noname001/hub/comm/tmphub/mdpservice"
)

type HubParams struct {
	RootContext   context.Context
	RootLogger    *logging.WrappedLogger
	RootLogPrefix string

	// TODO
	CfgRoot *rawconfig.ConfigRoot
}

type Hub_t struct {
	context context.Context
	cancel  context.CancelFunc

	logger    *logging.WrappedLogger
	logPrefix string

	commBroker   *comm.BrokerWrapper
	pubsubServer *comm.PubSubServer

	tmpHubProvider *tmphubMdp.TmpHubProvider
}

func NewHub(params *HubParams) (*Hub_t, error) {
	var err error

	hub := &Hub_t{}
	hub.context, hub.cancel = context.WithCancel(params.RootContext)
	hub.logger, hub.logPrefix = params.RootLogger, fmt.Sprintf("%s.hub", params.RootLogPrefix)

	hub.commBroker, err = comm.NewBrokerWrapper(comm.BrokerWrapperConfig{
		Context: hub.context,
		Verbose: params.CfgRoot.Hub.CommVerbose,
		BrokerHost: params.CfgRoot.Hub.BrokerHost,
	})
	if err != nil {
		hub.logger.Errorf("[%s] comm.NewBrokerWrapper err, %s", hub.logPrefix, err.Error())
		hub._abort()
		return nil, err
	}

	hub.pubsubServer, err = comm.NewPubsubServer(comm.PubSubServerConfig{
		Verbose: params.CfgRoot.Hub.CommVerbose,
		SnapshotHost: params.CfgRoot.Hub.SnapshotHost,
		PublisherHost: params.CfgRoot.Hub.PublisherHost,
		CollectorHost: params.CfgRoot.Hub.CollectorHost,
	})
	if err != nil {
		hub.logger.Errorf("[%s] comm.NewPubsubServer err, %s", hub.logPrefix, err.Error())
		hub._abort()
		return nil, err
	}

	hub.tmpHubProvider, err = tmphubMdp.NewTmpHubProvider(&tmphubMdp.TmpHubProviderParams{
		Context: hub.context,
		Logger: hub.logger, LogPrefix: hub.logPrefix,

		BrokerHost: params.CfgRoot.Hub.TmpHubProviderBrokerHost,
		DataHandler: &tmphubProvider_dataHandler_t{
			commBroker: hub.commBroker,
		},

		Verbose: params.CfgRoot.Hub.CommVerbose,
		RetryBackoff: tmphubMdp.DEFAULT_RETRY_BACKOFF,
	})
	if err != nil {
		hub.logger.Errorf("[%s] tmphubMdp.NewTmpHubProvider err, %s", hub.logPrefix, err.Error())
		hub._abort()
		return nil, err
	}

	hub.logger.Infof("[%s] initialized", hub.logPrefix)
	return hub, nil
}

func (hub *Hub_t) Start() (error) {
	var err error

	err = hub.commBroker.Start()
	if err != nil {
		hub.logger.Errorf("[%s] commBroker.Start err, %s", hub.logPrefix, err.Error())
		return err
	}

	hub.pubsubServer.Start()

	err = hub.tmpHubProvider.Connect()
	if err != nil {
		hub.logger.Errorf("[%s] tmpHubProvider.Connect err, %s", hub.logPrefix, err.Error())
		return err
	}

	hub.logger.Infof("[%s] started", hub.logPrefix)
	return nil
}

func (hub *Hub_t) Stop() {
	hub._cleanup()

	hub.logger.Infof("[%s] stopped", hub.logPrefix)
}

func (hub *Hub_t) _cleanup() {
	if hub.tmpHubProvider != nil { hub.tmpHubProvider.Disconnect() }

	if hub.pubsubServer != nil { hub.pubsubServer.Stop() }
	if hub.commBroker != nil { hub.commBroker.Stop() }

	hub.cancel()
}

func (hub *Hub_t) _abort() {
	hub._cleanup()

	hub.logger.Warningf("[%s] aborted", hub.logPrefix)
}
