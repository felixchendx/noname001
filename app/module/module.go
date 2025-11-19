package module

import (
	"context"

	"noname001/config/rawconfig"

	modDef "noname001/app/module/definition"
	"noname001/app/module/registry"
)

type ModuleRoot struct {
	modRegistry *registry.ModuleRegistry
}

func Initialize(ctx context.Context, cfg *rawconfig.ConfigRoot, runnerCfg *rawconfig.RunnerConfigRoot) (*ModuleRoot, error) {
	var err error

	modRoot := &ModuleRoot{}
	modRoot.modRegistry, err = registry.LoadModules(ctx, cfg, runnerCfg)
	if err != nil {
		return nil, err
	}

	return modRoot, nil
}

func (modRoot *ModuleRoot) States() (map[modDef.ModuleCode]string) {
	return modRoot.modRegistry.States()
}

// temp
func (modRoot *ModuleRoot) PlainStates() (map[string]string) {
	return modRoot.modRegistry.PlainStates()
}

func (modRoot *ModuleRoot) StartAll() (err error) {
	err = modRoot.modRegistry.StartAll()
	return
}

func (modRoot *ModuleRoot) StopAll() (err error) {
	err = modRoot.modRegistry.StopAll()
	return
}

// TODO: implement service locator
