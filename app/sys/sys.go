package sys

import (
	"context"

	"noname001/logging"
)

var (
	Bundle *Bundled
)

type Bundled struct {
	context context.Context
	cancel  context.CancelFunc

	store   *Store
	Service *SystemService
}

func Initialize(ctx context.Context) (error) {
	bundle := &Bundled{}
	bundle.context, bundle.cancel = context.WithCancel(ctx)

	store, err := NewStore(bundle.context, logging.Logger)
	if err != nil {
		return err
	}
	sev := store.DB.DBInit()
	if sev.IsError() {
		return sev.OriErr()
	}

	service, err := NewSystemService(bundle.context, store)
	if err != nil {
		return err
	}

	bundle.store = store
	bundle.Service = service

	Bundle = bundle

	sessionCleanupWorker(bundle.context)

	return nil
}

func (bundle *Bundled) Close() {
	bundle.Service.Stop()
	bundle.store.Close()
	bundle.cancel()
}
