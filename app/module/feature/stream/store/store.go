package store

import (
	"context"

	"noname001/logging"

	"noname001/app/module/feature/stream/store/sqlite"
)

type StoreParams struct {
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
}
type Store struct {
	context   context.Context
	cancel    context.CancelFunc

	logger    *logging.WrappedLogger
	logPrefix string

	DB        *sqlite.DB
}

func NewStore(params *StoreParams) (*Store, error) {
	store := &Store{}
	store.context, store.cancel = context.WithCancel(params.Context)
	store.logger = params.Logger
	store.logPrefix = params.LogPrefix + ".store"

	db, err := sqlite.NewDB(&sqlite.DBParams{
		Context: store.context,
		Logger: store.logger,
		LogPrefix: store.logPrefix,
	})
	if err != nil {
		return nil, err
	}

	store.DB = db

	return store, nil
}

func (store *Store) Open() (err error) {
	return
}

func (store *Store) Close() (err error) {
	store.DB.Close()
	store.cancel()
	return
}
