package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

func (db *DB) Wall__Insert(header *WallPE, items []*WallItemPE) (*DBEvent) {
	dbev := db.NewEvent("Wall__Insert")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	tx, err := db.Conn.BeginTx(qctx, db.DefaultTransactionOptions())
	if err != nil {
		db.LogError(dbev, err, "tx", "begin")
		return dbev
	}
	defer tx.Rollback()

	headerArgs := []any{
		sql.Named("id", header.ID),
		sql.Named("code", header.Code),
		sql.Named("name", header.Name),
		sql.Named("state", header.State),
		sql.Named("note", header.Note),

		sql.Named("wall_layout_id", header.WallLayoutID),
	}
	headerStmt := `
		INSERT INTO wall (
			id,
			code, name, state,
			note,

			wall_layout_id
		) VALUES (
			@id,
			@code, @name, @state,
			@note,

			@wall_layout_id
		);
	`
	result, err := tx.Exec(headerStmt, headerArgs...)
	if err != nil {
		db.LogError(dbev, err, headerStmt, headerArgs)
		return dbev
	}
	_ = result


	itemArgs := make([][]any, 0, len(items))
	for _, item := range items {
		itemArgs = append(itemArgs, []any{
			sql.Named("id", item.ID),
			sql.Named("wall_id", header.ID),
			sql.Named("idx", item.Index),
			sql.Named("source_node_id", item.SourceNodeID),
			sql.Named("stream_code", item.StreamCode),
		})
	}
	itemStmt := `
		INSERT INTO wall_item (
			id,
			wall_id, idx,

			source_node_id, stream_code
		) VALUES (
			@id,
			@wall_id, @idx,

			@source_node_id, @stream_code
		);
	`

	preppedItemStmt, err := tx.Prepare(itemStmt)
	if err != nil {
		db.LogError(dbev, err, "tx", "prep")
		return dbev
	}
	defer preppedItemStmt.Close()

	for _, itemArg := range itemArgs {
		_, err := preppedItemStmt.Exec(itemArg...)
		if err != nil {
			db.LogError(dbev, err, "tx", "prep.exec")
			return dbev
		}
	}

	if err := tx.Commit(); err != nil {
		db.LogError(dbev, err, "tx", "commit")
		return dbev
	}

	return dbev
}

func (db *DB) Wall__Get(headerID string, withItems bool) (*WallPE, []*WallItemPE, *DBEvent) {
	dbev := db.NewEvent("Wall__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	headerArgs := []any{headerID}
	headerStmt := `SELECT * FROM wall WHERE id = ?;`

	headerEntity := &WallPE{}
	if err := db.Conn.QueryRowContext(qctx, headerStmt, headerArgs...).Scan(
		headerEntity.fullScan()...
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, dbev
		} else {
			db.LogError(dbev, err, headerStmt, headerArgs)
			return nil, nil, dbev
		}
	}

	if !withItems {
		return headerEntity, nil, dbev
	}


	itemArgs := []any{headerEntity.ID}
	itemStmt := `
		SELECT *
		FROM wall_item
		WHERE wall_id = ?
		ORDER BY idx ASC
		;
	`

	// new ctx ?
	rows, err := db.Conn.QueryContext(qctx, itemStmt, itemArgs...)
	if err != nil {
		db.LogError(dbev, err, itemStmt, itemArgs)
		return nil, nil, dbev
	}

	itemEntities := make([]*WallItemPE, 0)
	for rows.Next() {
		var itemEntity *WallItemPE = &WallItemPE{}
		if err := rows.Scan(
			itemEntity.fullScan()...
		); err != nil {
			db.LogError(dbev, err, "item", "rows.scan")
			return nil, nil, dbev
		}
		itemEntities = append(itemEntities, itemEntity)
	}

	if err := rows.Err(); err != nil {
		db.LogError(dbev, err, "item", "rows.err")
		return nil, nil, dbev
	}

	return headerEntity, itemEntities, dbev
}

func (db *DB) Wall__GetByCode(headerCode string, withItems bool) (*WallPE, []*WallItemPE, *DBEvent) {
	dbev := db.NewEvent("Wall__GetByCode")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	headerArgs := []any{headerCode}
	headerStmt := `SELECT * FROM wall WHERE code = ?;`

	headerEntity := &WallPE{}
	if err := db.Conn.QueryRowContext(qctx, headerStmt, headerArgs...).Scan(
		headerEntity.fullScan()...
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, dbev
		} else {
			db.LogError(dbev, err, headerStmt, headerArgs)
			return nil, nil, dbev
		}
	}

	if !withItems {
		return headerEntity, nil, dbev
	}


	itemArgs := []any{headerEntity.ID}
	itemStmt := `
		SELECT *
		FROM wall_item
		WHERE wall_id = ?
		ORDER BY idx ASC
		;
	`

	// new ctx ?
	rows, err := db.Conn.QueryContext(qctx, itemStmt, itemArgs...)
	if err != nil {
		db.LogError(dbev, err, itemStmt, itemArgs)
		return nil, nil, dbev
	}

	itemEntities := make([]*WallItemPE, 0)
	for rows.Next() {
		var itemEntity *WallItemPE = &WallItemPE{}
		if err := rows.Scan(
			itemEntity.fullScan()...
		); err != nil {
			db.LogError(dbev, err, "item", "rows.scan")
			return nil, nil, dbev
		}
		itemEntities = append(itemEntities, itemEntity)
	}

	if err := rows.Err(); err != nil {
		db.LogError(dbev, err, "item", "rows.err")
		return nil, nil, dbev
	}

	return headerEntity, itemEntities, dbev
}

func (db *DB) Wall__Update(header *WallPE, items []*WallItemPE) (*DBEvent) {
	dbev := db.NewEvent("Wall__Update")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	tx, err := db.Conn.BeginTx(qctx, db.DefaultTransactionOptions())
	if err != nil {
		db.LogError(dbev, err, "tx", "begin")
		return dbev
	}
	defer tx.Rollback()

	headerArgs := []any{
		sql.Named("code", header.Code),
		sql.Named("name", header.Name),
		sql.Named("state", header.State),
		sql.Named("note", header.Note),

		sql.Named("wall_layout_id", header.WallLayoutID),
	}
	headerArgs = append(headerArgs, sql.Named("id", header.ID))
	headerStmt := `
		UPDATE wall SET
			code = @code, name = @name, state = @state,
			note = @note,

			wall_layout_id = @wall_layout_id,

			updated_ts = CURRENT_TIMESTAMP
		WHERE id = @id
		;
	`

	result, err := tx.Exec(headerStmt, headerArgs...)
	if err != nil {
		db.LogError(dbev, err, headerStmt, headerArgs)
		return dbev
	}
	_ = result


	deleteAllItemArgs := []any{&header.ID}
	deleteAllItemStmt := `DELETE FROM wall_item WHERE wall_id = ?;`
	result2, err := tx.Exec(deleteAllItemStmt, deleteAllItemArgs...)
	if err != nil {
		db.LogError(dbev, err, deleteAllItemStmt, deleteAllItemArgs)
		return dbev
	}
	_ = result2


	itemArgs := make([][]any, 0, len(items))
	for _, item := range items {
		itemArgs = append(itemArgs, []any{
			sql.Named("id", item.ID),
			sql.Named("wall_id", header.ID),
			sql.Named("idx", item.Index),
			sql.Named("source_node_id", item.SourceNodeID),
			sql.Named("stream_code", item.StreamCode),
		})
	}
	itemStmt := `
		INSERT INTO wall_item (
			id,
			wall_id, idx,

			source_node_id, stream_code
		) VALUES (
			@id,
			@wall_id, @idx,

			@source_node_id, @stream_code
		);
	`

	preppedItemStmt, err := tx.Prepare(itemStmt)
	if err != nil {
		db.LogError(dbev, err, "tx", "prep")
		return dbev
	}
	defer preppedItemStmt.Close()

	for _, itemArg := range itemArgs {
		_, err := preppedItemStmt.Exec(itemArg...)
		if err != nil {
			db.LogError(dbev, err, "tx", "prep.exec")
			return dbev
		}
	}

	if err := tx.Commit(); err != nil {
		db.LogError(dbev, err, "tx", "commit")
		return dbev
	}

	return dbev
}

func (db *DB) Wall__Delete(headerID string) (*DBEvent) {
	dbev := db.NewEvent("Wall__Delete")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	tx, err := db.Conn.BeginTx(qctx, db.DefaultTransactionOptions())
	if err != nil {
		db.LogError(dbev, err, "tx", "begin")
		return dbev
	}
	defer tx.Rollback()

	args := []any{sql.Named("wall_id", headerID)}
	stmt := `
		DELETE FROM wall_item WHERE wall_id = @wall_id;
		DELETE FROM wall WHERE id = @wall_id;
	`

	result, err := tx.Exec(stmt, args...)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return dbev
	}
	_ = result

	if err := tx.Commit(); err != nil {
		db.LogError(dbev, err, "tx", "commit")
		return dbev
	}

	return dbev
}


func (db *DB) WallItem__Get(id string) (*WallItemPE, *DBEvent) {
	dbev := db.NewEvent("WallItem__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{id}
	stmt := `SELECT * FROM wall_item WHERE id = ?;`

	retEntity := &WallItemPE{}
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

func (db *DB) WallItem__Update(id string, entity *WallItemPE) (*WallItemPE, *DBEvent) {
	dbev := db.NewEvent("WallItem__Update")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.SourceNodeID,
		&entity.StreamCode,
	}
	args = append(args, id)
	stmt := `
		UPDATE wall_item SET
			source_node_id = ?,
			stream_code = ?,

			updated_ts = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING *
		;
	`

	retEntity := &WallItemPE{}
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
