package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

func (db *DB) StreamItem__AtomicInsert(entity *StreamItemPE) (*StreamItemPE, *DBEvent) {
	dbev := db.NewEvent("StreamItem__AtomicInsert")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.ID,
		&entity.StreamGroupID,
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.SourceType,
		&entity.DeviceCode, &entity.DeviceChannelID, &entity.DeviceStreamType,
		&entity.ExternalURL,
		&entity.Filepath,
		&entity.EmbeddedFilepath,
	}
	stmt := `
		INSERT INTO stream_item (
				id,
				stream_group_id,
				code, name, state,
				note,

				source_type,
				device_code, device_channel_id, device_stream_type,
				external_url,
				filepath,
				embedded_filepath
			) VALUES (
				?,
				?,
				?, ?, ?,
				?,

				?,
				?, ?, ?,
				?,
				?,
				?
			)
		RETURNING *
		;
	`
	
	retEntity := &StreamItemPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamItem__EntityFullScan(retEntity)...
	); err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	return retEntity, dbev
}

func (db *DB) StreamItem__Get(id string) (*StreamItemPE, *DBEvent) {
	dbev := db.NewEvent("StreamItem__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{id}
	stmt := `SELECT * FROM stream_item WHERE id = ?;`

	retEntity := &StreamItemPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamItem__EntityFullScan(retEntity)...
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

func (db *DB) StreamItem__GetByCode(code string) (*StreamItemPE, *DBEvent) {
	dbev := db.NewEvent("StreamItem__GetByCode")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{code}
	stmt := `SELECT * FROM stream_item WHERE code = ?;`

	retEntity := &StreamItemPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamItem__EntityFullScan(retEntity)...
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

func (db *DB) StreamItem__GetByStreamGroupID(sgID string) ([]*StreamItemPE, *DBEvent) {
	dbev := db.NewEvent("StreamItem__GetByStreamGroupID")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{sgID}
	stmt := `SELECT * FROM stream_item WHERE stream_group_id = ?;`

	rows, err := db.Conn.QueryContext(qctx, stmt, args...)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	entities := make([]*StreamItemPE, 0)
	for rows.Next() {
		var entity *StreamItemPE = &StreamItemPE{}
		if err := rows.Scan(
			db.StreamItem__EntityFullScan(entity)...
		); err != nil {
			db.LogError(dbev, err, "mainQ", "rows.scan")
			return nil, dbev
		}

		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		db.LogError(dbev, err, "mainQ", "rows.err")
		return nil, dbev
	}

	return entities, dbev
}

func (db *DB) StreamItem__AtomicUpdate(entity *StreamItemPE) (*StreamItemPE, *DBEvent) {
	dbev := db.NewEvent("StreamItem__AtomicUpdate")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.StreamGroupID,
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.SourceType,
		&entity.DeviceCode, &entity.DeviceChannelID, &entity.DeviceStreamType,
		&entity.ExternalURL,
		&entity.Filepath,
		&entity.EmbeddedFilepath,
	}
	args = append(args, &entity.ID)
	stmt := `
		UPDATE stream_item SET
			stream_group_id = ?,
			code = ?, name = ?, state = ?,
			note = ?,

			source_type = ?,
			device_code = ?, device_channel_id = ?, device_stream_type = ?,
			external_url = ?,
			filepath = ?,
			embedded_filepath = ?,

			updated_ts = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING *
		;
	`

	retEntity := &StreamItemPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamItem__EntityFullScan(retEntity)...
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

func (db *DB) StreamItem__AtomicDelete(id string) (*DBEvent) {
	dbev := db.NewEvent("StreamItem__AtomicDelete")

	args := []any{id}
	stmt := `DELETE FROM stream_item WHERE id = ?;`

	err := db.ExecuteWithArgs(stmt, args)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return dbev
	}

	return dbev
}
