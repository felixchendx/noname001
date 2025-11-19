package sqlite

import (
	"context"
	"path/filepath"

	"noname001/filesystem"
	"noname001/logging"

	baseStore "noname001/app/base/store"
)

const (
	MODULE_CODE  = "stream"
	DB_FILENAME  = "stream.dat"

	QUERY_TIMEOUT = baseStore.SQLITE_QUERY_TIMEOUT
)

type DBParams struct {
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
}
type DB struct {
	*baseStore.SQLiteDB
}

func NewDB(params *DBParams) (*DB, error) {
	sqliteDBFileLocation := filepath.Join(filesystem.DBDir, DB_FILENAME)

	sqliteDB, err := baseStore.NewSQLiteDB(baseStore.SQLiteDBParams{
		Context: params.Context,
		Logger: params.Logger,
		LogPrefix: params.LogPrefix + ".db",

		SQLiteDBFileLocation: sqliteDBFileLocation,
	})
	if err != nil {
		return nil, err
	}

	return &DB{sqliteDB}, nil
}
