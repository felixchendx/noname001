package app

import (
	"context"
	"fmt"

	"noname001/logging"
	"noname001/config/rawconfig"

	appSys "noname001/app/sys"
	appMod "noname001/app/module"
)

type AppParams struct {
	RootContext   context.Context
	RootLogger    *logging.WrappedLogger
	RootLogPrefix string

	// TODO
	CfgRoot       *rawconfig.ConfigRoot
	RunnerCfgRoot *rawconfig.RunnerConfigRoot
}

type App_t struct {
	context context.Context
	cancel  context.CancelFunc

	logger    *logging.WrappedLogger
	logPrefix string

	modRoot *appMod.ModuleRoot
}

func NewApp(params *AppParams) (*App_t, error) {
	var err error

	app := &App_t{}
	app.context, app.cancel = context.WithCancel(params.RootContext)
	app.logger, app.logPrefix = params.RootLogger, fmt.Sprintf("%s.app", params.RootLogPrefix)

	err = appSys.Initialize(app.context)
	if err != nil {
		app.logger.Errorf("[%s] appSys.Initialize err, %s", app.logPrefix, err.Error())
		app._abort()
		return nil, err
	}

	app.modRoot, err = appMod.Initialize(app.context, params.CfgRoot, params.RunnerCfgRoot)
	if err != nil {
		app.logger.Errorf("[%s] appMod.Initialize err, %s", app.logPrefix, err.Error())
		app._abort()
		return nil, err
	}

	app.logger.Infof("[%s] initialized", app.logPrefix)
	return app, nil
}

func (app *App_t) Start() (error) {
	var err error

	err = app.modRoot.StartAll()
	if err != nil {
		app.logger.Errorf("[%s] modRoot.StartAll err, %s", app.logPrefix, err.Error())
		return err
	}

	app.logger.Infof("[%s] started", app.logPrefix)
	return nil
}

func (app *App_t) Stop() {
	app._cleanup()

	app.logger.Infof("[%s] stopped", app.logPrefix)
}

func (app *App_t) ModStates() (map[string]string) {
	return app.modRoot.PlainStates()
}

func (app *App_t) _cleanup() {
	if app.modRoot != nil { app.modRoot.StopAll() }
	if appSys.Bundle != nil { appSys.Bundle.Close() }

	app.cancel()
}

func (app *App_t) _abort() {
	app._cleanup()

	app.logger.Warningf("[%s] aborted", app.logPrefix)
}
