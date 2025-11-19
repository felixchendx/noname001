package node

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"

	"noname001/config/rawconfig"
	"noname001/logging"

	nodeTyping "noname001/node/base/typing"
	nodeEv     "noname001/node/event"

	"noname001/hub"
	"noname001/app"
)

type NodeParams struct {
	RootContext   context.Context
	RootLogger    *logging.WrappedLogger
	RootLogPrefix string

	CfgRoot *rawconfig.ConfigRoot
}

type Node_t struct {
	context context.Context
	cancel  context.CancelFunc

	logger    *logging.WrappedLogger
	logPrefix string

	cfgRoot *rawconfig.ConfigRoot

	cron     *cron.Cron
	cronJobs map[string]cron.EntryID

	id    string
	name  string
	state nodeTyping.NodeState

	ips                   []string
	ipCollectionHistories []*ipCollectionHistory_t

	evHub *nodeEv.EventHub

	commBundle *nodeCommBundle_t

	hub *hub.Hub_t
	app *app.App_t

	snapshot *nodeTyping.BaseNodeSnapshot
}

func NewNode(params *NodeParams) (*Node_t, error) {
	var (
		err  error
		node *Node_t
	)
	
	node = &Node_t{}
	node.context, node.cancel = context.WithCancel(params.RootContext)
	node.logger, node.logPrefix = params.RootLogger, fmt.Sprintf("%s.node", params.RootLogPrefix)

	node.cfgRoot = params.CfgRoot

	// TODO: node time object, time related source of truth
	node.cron = cron.New(
		cron.WithLocation(params.CfgRoot.Global.TimeLoc),
		cron.WithSeconds(),
	)
	node.cronJobs = make(map[string]cron.EntryID)

	node.id, node.name = node.cfgRoot.Node.ID, node.cfgRoot.Node.Name

	node.initIPWatcher()

	node.commBundle, err = node.initComm()
	if err != nil {
		node.logger.Errorf("[%s] node.initComm err, %s", node.logPrefix, err.Error())
		node._abort()
		return nil, err
	}

	node.evHub = nodeEv.NewEventHub(&nodeEv.EventHubParams{
		ParentContext: node.context,
		Logger: node.logger, LogPrefix: node.logPrefix,
	})

	node.injectCommConf(params.CfgRoot)

	node.state = nodeTyping.NODE_STATE__INIT
	node.logger.Infof("[%s] initialized", node.logPrefix)
	return node, nil
}

func (node *Node_t) Start() (error) {
	var err error

	node.evHub.Open()

	err = node.commBundle.connect()
	if err != nil {
		node.logger.Errorf("[%s] commBundle.connect err, %s", node.logPrefix, err.Error())
		return err
	}

	node.startIPWatcher()
	node.cron.Start()

	node.state = nodeTyping.NODE_STATE__START
	node.logger.Infof("[%s] started", node.logPrefix)
	return nil
}

func (node *Node_t) Stop() {
	node._cleanup()

	node.state = nodeTyping.NODE_STATE__STOP
	node.logger.Infof("[%s] stopped",  node.logPrefix)
}

func (node *Node_t) Ready() {
	node._announce(nodeTyping.NODE_EVENT_CODE__READY)
	node.setupHeartbeat()

	node.state = nodeTyping.NODE_STATE__READY
	node.logger.Infof("[%s] NODE '%s' READY!", node.logPrefix, node.id)
}

func (node *Node_t) Shutdown() {
	node._announce(nodeTyping.NODE_EVENT_CODE__SHUTDOWN)

	node.state = nodeTyping.NODE_STATE__SHUTDOWN
	node.logger.Infof("[%s] Node '%s' SHUTDOWN!", node.logPrefix, node.id)
}

func (node *Node_t) InjectHubInstance(hub *hub.Hub_t) {
	node.hub = hub
}
func (node *Node_t) InjectAppInstance(app *app.App_t) {
	node.app = app
}

func (node *Node_t) _cleanup() {
	node.cron.Stop()

	if node.evHub != nil { node.evHub.Close() }
	if node.commBundle != nil { node.commBundle.disconnect() }

	node.cancel()
}

func (node *Node_t) _abort() {
	node._cleanup()

	node.state = nodeTyping.NODE_STATE__ABORT
	node.logger.Warningf("[%s] aborted",  node.logPrefix)
}
