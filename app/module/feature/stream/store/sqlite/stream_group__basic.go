package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

func (db *DB) StreamGroup__AtomicInsert(entity *StreamGroupPE) (*StreamGroupPE, *DBEvent) {
	dbev := db.NewEvent("StreamGroup__AtomicInsert")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.ID,
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.StreamProfileID,
	}
	stmt := `
		INSERT INTO stream_group (
				id,
				code, name, state,
				note,

				stream_profile_id
			) VALUES (
				?,
				?, ?, ?,
				?,

				?
			)
		RETURNING *
		;
	`
	
	retEntity := &StreamGroupPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamGroup__EntityFullScan(retEntity)...
	); err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	return retEntity, dbev
}

func (db *DB) StreamGroup__Get(id string) (*StreamGroupPE, *DBEvent) {
	dbev := db.NewEvent("StreamGroup__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{id}
	stmt := `SELECT * FROM stream_group WHERE id = ?;`

	retEntity := &StreamGroupPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamGroup__EntityFullScan(retEntity)...
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

func (db *DB) StreamGroup__GetByCode(code string) (*StreamGroupPE, *DBEvent) {
	dbev := db.NewEvent("StreamGroup__GetByCode")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{code}
	stmt := `SELECT * FROM stream_group WHERE code = ?;`

	retEntity := &StreamGroupPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamGroup__EntityFullScan(retEntity)...
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

func (db *DB) StreamGroup__AtomicUpdate(entity *StreamGroupPE) (*StreamGroupPE, *DBEvent) {
	dbev := db.NewEvent("StreamGroup__AtomicUpdate")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.StreamProfileID,
	}
	args = append(args, &entity.ID)
	stmt := `
		UPDATE stream_group SET
			code = ?, name = ?, state = ?,
			note = ?,

			stream_profile_id = ?,

			updated_ts = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING *
		;
	`

	retEntity := &StreamGroupPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamGroup__EntityFullScan(retEntity)...
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
