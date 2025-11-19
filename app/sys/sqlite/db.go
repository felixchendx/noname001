package sqlite

import (
	"context"
	"path/filepath"

	"noname001/filesystem"
	"noname001/logging"

	baseStore "noname001/app/base/store"
)

const (
	MODULE_CODE  = "sys"
	DB_FILENAME  = "sys.dat"

	QUERY_TIMEOUT = baseStore.SQLITE_QUERY_TIMEOUT
)

// === debe ===
type DB struct {
	*baseStore.SQLiteDB
}

func NewDB(ctx context.Context, logger *logging.WrappedLogger) (*DB, error) {
	sqliteDBFileLocation := filepath.Join(filesystem.DBDir, DB_FILENAME)

	sqliteDB, err := baseStore.NewSQLiteDB(baseStore.SQLiteDBParams{
		Context: ctx,
		Logger: logger,
		LogPrefix: "app.sys.db",

		SQLiteDBFileLocation: sqliteDBFileLocation,
	})
	if err != nil {
		return nil, err
	}

	return &DB{sqliteDB}, nil
}
