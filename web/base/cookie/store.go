package cookie

import (
	// "encoding/gob"

	"noname001/logging"
)

type CookieStoreParams struct {
	Logger *logging.WrappedLogger
	LogPrefix string
}
type CookieStore struct {
	logger    *logging.WrappedLogger
	logPrefix string
}

func NewCookieStore(params *CookieStoreParams) (*CookieStore) {
	// gob.Register(&sessionCookieValue{})

	store := &CookieStore{}
	store.logger = params.Logger
	store.logPrefix = params.LogPrefix + ".cookie"

	return store
}
