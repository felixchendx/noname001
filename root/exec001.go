package root

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"noname001/config/rawconfig"
	"noname001/dilemma"
	localEv "noname001/dilemma/event"
	"noname001/filesystem"
	"noname001/logging"

	"noname001/app"
	"noname001/hub"
	"noname001/node"
	"noname001/web"
)

var (
	exec001 *exec001_t
)

type exec001_t struct {
	context context.Context
	cancel  context.CancelFunc

	logger    *logging.WrappedLogger
	logPrefix string

	timezone  *time.Location
	startTime time.Time

	evHub *localEv.EventHub

	node *node.Node_t
	hub  *hub.Hub_t
	app  *app.App_t

	web *web.WebService
}

func Start(cfgRoot *rawconfig.ConfigRoot, runnerCfgRoot *rawconfig.RunnerConfigRoot) error {
	var err error

	err = _init(cfgRoot, runnerCfgRoot)
	if err != nil {
		return err
	}

	err = _start()
	if err != nil {
		return err
	}

	return nil
}

func _init(cfgRoot *rawconfig.ConfigRoot, runnerCfgRoot *rawconfig.RunnerConfigRoot) error {
	var err error

	filesystem.SetRootDir(strings.TrimSpace(cfgRoot.Global.RootDirectory))
	err = filesystem.PrepareAll()
	if err != nil {
		// TODO
		_stop()
		return err
	}

	logging.ConfigureLogging(&logging.LoggingConfig{
		LogTo:    cfgRoot.Logging.LogTo,
		LogLevel: cfgRoot.Logging.LogLevel,
		LogTzStr: cfgRoot.Logging.LogTimezone,
	})

	dilemma.InjectTemporaryLogger(logging.Logger) // temp

	exec001 = &exec001_t{}
	exec001.context, exec001.cancel = context.WithCancel(context.Background())
	exec001.logger, exec001.logPrefix = logging.Logger, "exec001"

	exec001.logger.Infof("[%s] Initializing... with root dir: '%s'", exec001.logPrefix, cfgRoot.Global.RootDirectory)
	// TODO, print necessary info in list format

	// TODO: to be replaced with sanitized config, and then move validation there
	exec001.timezone, err = time.LoadLocation(cfgRoot.Global.Timezone)
	if err != nil {
		exec001.logger.Errorf("[%s] Invalid global timezone '%s'", exec001.logPrefix, cfgRoot.Global.Timezone)
		_stop()
		return err
	}
	exec001.startTime = time.Now().UTC()

	// TODO: streamlined when replacing this with sanitized config
	cfgRoot.Global.TimeLoc = exec001.timezone

	// === evHub - begin ===
	exec001.evHub = localEv.NewEventHub(&localEv.EventHubParams{
		ParentContext: exec001.context,
		Logger:        exec001.logger, LogPrefix: exec001.logPrefix,
	})
	// === evHub - end =====

	// === node - begin ===
	exec001.node, err = node.NewNode(&node.NodeParams{
		RootContext:   exec001.context,
		RootLogger:    exec001.logger,
		RootLogPrefix: exec001.logPrefix,

		// TODO: sanitized config
		CfgRoot: cfgRoot,
	})
	if err != nil {
		_stop()
		return err
	}
	// === node - end =====

	// === hub - begin ===
	if cfgRoot.Hub.Enabled {
		exec001.hub, err = hub.NewHub(&hub.HubParams{
			RootContext:   exec001.context,
			RootLogger:    exec001.logger,
			RootLogPrefix: exec001.logPrefix,

			// TODO: sanitized config
			CfgRoot: cfgRoot,
		})
		if err != nil {
			_stop()
			return err
		}

		exec001.node.InjectHubInstance(exec001.hub)
	}
	// === hub - end =====

	// === app - begin ===
	if cfgRoot.Application.Enabled {
		exec001.app, err = app.NewApp(&app.AppParams{
			RootContext:   exec001.context,
			RootLogger:    exec001.logger,
			RootLogPrefix: exec001.logPrefix,

			// TODO: sanitized config
			CfgRoot:       cfgRoot,
			RunnerCfgRoot: runnerCfgRoot,
		})
		if err != nil {
			_stop()
			return err
		}

		exec001.node.InjectAppInstance(exec001.app)
	}
	// === app - end =====

	// === web - begin ===
	// TODO: start web first, for start status page
	//       but need to sweep stuffs that call other un-initialized parts
	if cfgRoot.Web.Enabled {
		exec001.web, err = web.NewWebService(&web.WebServiceParams{
			RootContext:   exec001.context,
			RootLogger:    exec001.logger, // TODO: independent logger for web stuffs
			RootLogPrefix: exec001.logPrefix,

			Name:       "noname",
			ListenHost: fmt.Sprintf("%s:%s", cfgRoot.Web.Hostname, cfgRoot.Web.Port),

			Behind7Proxies:  cfgRoot.Web.Behind7Proxies,
			ServerStartTime: exec001.startTime,

			TempModStates: make(map[string]string),
		})
		if err != nil {
			_stop()
			return err
		}
	}
	// === web - end =====

	exec001.logger.Infof("[%s] Initialized", exec001.logPrefix)
	return nil
}

func _start() error {
	exec001.logger.Infof("[%s] Starting...", exec001.logPrefix)

	var err error

	err = exec001.evHub.Open()
	if err != nil {
		_stop()
		return err
	}

	err = exec001.node.Start()
	if err != nil {
		_stop()
		return err
	}

	if exec001.hub != nil {
		err = exec001.hub.Start()
		if err != nil {
			_stop()
			return err
		}
	}

	if exec001.app != nil {
		err = exec001.app.Start()
		if err != nil {
			_stop()
			return err
		}
	}

	if exec001.web != nil {
		// TODO: proper state signaling
		modStates := make(map[string]string)
		if exec001.app != nil {
			modStates = exec001.app.ModStates()
		}
		exec001.web.InjectTempModStates(modStates)

		// TODO: check web server properly started, and _stop() if necessary
		exec001.web.Serve()
	}

	exec001.node.Ready()
	exec001.logger.Infof("[%s] Start OK! To exit, send SIGINT (CTRL+C) or SIGTERM", exec001.logPrefix)

	exitSigChan := make(chan os.Signal, 1)
	signal.Notify(exitSigChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

rootLoop:
	for {
		select {
		case exitSig := <-exitSigChan:
			exec001.node.Shutdown() // TODO: node event wrong order
			logging.Logger.Infof("[%s] Received signal '%s'. Initiating graceful shutdown...", exec001.logPrefix, exitSig)
			_stop()
			break rootLoop
		}
	}

	// TODO: shutdownLoop

	exec001.logger.Infof("[%s] Exit normal", exec001.logPrefix)
	return nil
}

func _stop() {
	exec001.logger.Infof("[%s] Stopping...", exec001.logPrefix)

	if exec001.web != nil {
		exec001.web.Shutdown()
	}
	if exec001.app != nil {
		exec001.app.Stop()
	}
	if exec001.hub != nil {
		exec001.hub.Stop()
	}
	if exec001.node != nil {
		exec001.node.Stop()
	}
	if exec001.evHub != nil {
		exec001.evHub.Close()
	}

	exec001.cancel()

	exec001.logger.Infof("[%s] Stopped", exec001.logPrefix)
}
