package sys

import (
	"context"

	"noname001/logging"

	"noname001/app/sys/sqlite"
)

type Store struct {
	context context.Context
	cancel  context.CancelFunc

	logger *logging.WrappedLogger

	DB      *sqlite.DB
}

func NewStore(ctx context.Context, logger *logging.WrappedLogger) (*Store, error) {
	store := &Store{}
	store.context, store.cancel = context.WithCancel(ctx)
	store.logger = logger
	
	db, err := sqlite.NewDB(store.context, store.logger)
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
