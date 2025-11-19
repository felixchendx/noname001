package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

type PingPE struct {
	Reply     string
	Timestamp string
}

func (db *DB) Ping() (*PingPE, *DBEvent) {
	dbev := db.NewEvent("Ping")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{"pong"}
	stmt := `SELECT id, CURRENT_TIMESTAMP FROM ping WHERE id = ?;`

	var entity *PingPE = &PingPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		&entity.Reply,
		&entity.Timestamp,
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
