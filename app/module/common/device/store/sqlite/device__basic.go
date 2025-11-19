package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

func (db *DB) Device__AtomicInsert(entity *DevicePE) (*DevicePE, *DBEvent) {
	dbev := db.NewEvent("Device__AtomicInsert")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.ID,
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.Protocol,
		&entity.Hostname, &entity.Port,
		&entity.Username, &entity.Password,
		&entity.Brand,

		&entity.FallbackRTSPPort,
	}
	stmt := `
		INSERT INTO device (
				id,
				code, name, state,
				note,

				protocol,
				hostname, port,
				username, password,
				brand,

				fallbackRTSPPort
			) VALUES (
				?,
				?, ?, ?,
				?,

				?,
				?, ?,
				?, ?,
				?,

				?
			)
		RETURNING *
		;
	`
	
	retEntity := &DevicePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
	); err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	return retEntity, dbev
}

func (db *DB) Device__Get(id string) (*DevicePE, *DBEvent) {
	dbev := db.NewEvent("Device__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{id}
	stmt := `SELECT * FROM device WHERE id = ?;`

	retEntity := &DevicePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbev
		} else {
			db.LogError(dbev, err, stmt, args)
			return nil, dbev
		}
	}

	return retEntity, dbev
}

func (db *DB) Device__GetByCode(code string) (*DevicePE, *DBEvent) {
	dbev := db.NewEvent("Device__GetByCode")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{code}
	stmt := `SELECT * FROM device WHERE code = ?;`

	retEntity := &DevicePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbev
		} else {
			db.LogError(dbev, err, stmt, args)
			return nil, dbev
		}
	}

	return retEntity, dbev
}

func (db *DB) Device__AtomicUpdate(entity *DevicePE) (*DevicePE, *DBEvent) {
	dbev := db.NewEvent("Device__AtomicUpdate")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.Protocol,
		&entity.Hostname, &entity.Port,
		&entity.Username, &entity.Password,
		&entity.Brand,

		&entity.FallbackRTSPPort,
	}
	args = append(args, &entity.ID)
	stmt := `
		UPDATE device SET
			code = ?, name = ?, state = ?,
			note = ?,

			protocol = ?,
			hostname = ?, port = ?,
			username = ?, password = ?,
			brand = ?,

			fallbackRTSPPort = ?,

			updated_ts = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING *
		;
	`

	retEntity := &DevicePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbev
		} else {
			db.LogError(dbev, err, stmt, args)
			return nil, dbev
		}
	}

	return retEntity, dbev
}

func (db *DB) Device__AtomicDelete(id string) (*DBEvent) {
	dbev := db.NewEvent("Device__AtomicDelete")

	args := []any{id}
	stmt := `DELETE FROM device WHERE id = ?;`

	err := db.ExecuteWithArgs(stmt, args)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return dbev
	}

	return dbev
}
