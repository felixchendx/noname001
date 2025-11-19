package sqlite

import (
	"context"
	"database/sql"
)

func (db *DB) StreamGroup__CascadingDelete(headerID string) (*DBEvent) {
	dbev := db.NewEvent("StreamGroup__CascadingDelete")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	tx, err := db.Conn.BeginTx(qctx, db.DefaultTransactionOptions())
	if err != nil {
		db.LogError(dbev, err, "tx", "begin")
		return dbev
	}
	defer tx.Rollback()

	args := []any{sql.Named("stream_group_id", headerID)}
	stmt := `
		DELETE FROM stream_item WHERE stream_group_id = @stream_group_id;
		DELETE FROM stream_group WHERE id = @stream_group_id;
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
