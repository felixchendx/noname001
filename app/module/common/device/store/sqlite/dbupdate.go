package sqlite

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	dbUpdateErr001 = errors.New("missing critical information to do dbupdate.")
)

func (db *DB) DBUpdate() (*DBEvent) {
	dbev := db.NewEvent("DBUpdate")

	dbinfoPE, _dbev := db.getDBInfo()
	if _dbev.IsError() {
		return _dbev
	}
	if dbinfoPE == nil {
		db.LogError(dbev, dbUpdateErr001, "dbUpdateErr001", "")
		return dbev
	}

	currDBVersion := dbinfoPE.DBVersion
	var updateVersion int64 = 0


	updateVersion = 1
	if err := db.doDBUpdate(currDBVersion, updateVersion, `
		CREATE TABLE device
		(
			id         TEXT NOT NULL PRIMARY KEY,
			code       TEXT NOT NULL DEFAULT '',
			name       TEXT NOT NULL DEFAULT '',
			state      TEXT NOT NULL DEFAULT 'active',		-- 'active' | 'inactive'
			note       TEXT NOT NULL DEFAULT '',

			protocol   TEXT NOT NULL DEFAULT '',
			hostname   TEXT NOT NULL DEFAULT '',
			port       TEXT NOT NULL DEFAULT '',
			username   TEXT NOT NULL DEFAULT '',
			password   TEXT NOT NULL DEFAULT '',
			brand      TEXT NOT NULL DEFAULT '',

			fallbackRTSPPort TEXT NOT NULL DEFAULT '',

			created_ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			UNIQUE(code)
		);
		CREATE INDEX idx__device__state ON device(state);
		CREATE INDEX idx__device__brand ON device(brand);
	`); err != nil {
		db.LogError(dbev, err, "", updateVersion)
		return dbev
	}

	return dbev
}

func (db *DB) doDBUpdate(currDBVersion int64, dbUpdateVersion int64, sqlStatement string) (error) {
	dbev := db.NewEvent("doDBUpdate")
	
	if currDBVersion >= dbUpdateVersion {
		db.LogDebug(dbev, "[%04v] skipping... (%v >= %v)", dbUpdateVersion, currDBVersion, dbUpdateVersion)
		return nil
	}
	db.LogInfo(dbev, "[%04v] updating...", dbUpdateVersion)

	sqlStatement += fmt.Sprintf(
		`
		UPDATE dbinfo SET db_version = %s WHERE code = '%s';
		`,
		strconv.FormatInt(dbUpdateVersion, 10),
		MODULE_CODE,
	)

	err := db.Execute(sqlStatement)

	return err
}
