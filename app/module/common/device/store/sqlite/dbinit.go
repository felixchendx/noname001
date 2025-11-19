package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

type dbinfoPE struct {
	ID          string
	Code        string
	Name        string
	Type        string
	Version     string
	DBVersion   int64

	CreatedDttm string
	UpdatedDttm string
}
func (db *DB) getDBInfo() (*dbinfoPE, *DBEvent) {
	dbev := db.NewEvent("getDBInfo")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{MODULE_CODE}
	stmt := `SELECT * FROM dbinfo WHERE code = ?;`

	var entity *dbinfoPE = &dbinfoPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		&entity.ID,
		&entity.Code, &entity.Name, &entity.Type,
		&entity.Version, &entity.DBVersion,

		&entity.CreatedDttm, &entity.UpdatedDttm,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbev
		} else {
			db.LogError(dbev, err, stmt, args)
			return nil, dbev
		}
	}

	return entity, dbev
}

func (db *DB) DBInit() (*DBEvent) {
	dbev := db.NewEvent("DBInit")

	args := []any{
		sql.Named("id", "a45921e8-af1c-4771-a431-5eec870bc399"),
		sql.Named("mod_code", MODULE_CODE),
	}

	stmt := `
		CREATE TABLE IF NOT EXISTS dbinfo (
			id         TEXT NOT NULL PRIMARY KEY,

			code       TEXT NOT NULL DEFAULT '',
			name       TEXT NOT NULL DEFAULT '',
			type       TEXT NOT NULL DEFAULT '',		-- 'app' | 'mod'
			version    TEXT NOT NULL DEFAULT '',
			db_version TEXT NOT NULL DEFAULT '',

			created_ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		INSERT INTO dbinfo (id, code, name, type, version, db_version)
			SELECT @id, @mod_code, @mod_code, 'mod', '0.0.0', 0
		WHERE NOT EXISTS(
			SELECT 1 FROM dbinfo WHERE id = @id
		);
	`
	if err := db.ExecuteWithArgs(stmt, args); err != nil {
		db.LogError(dbev, err, stmt, args)
		return dbev
	}

	return dbev
}
