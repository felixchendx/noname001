package sqlite

import (
	"context"
)

func (db *DB) Coord__postStart_loadLiveStream() ([]*StreamItemPE, *DBEvent) {
	dbev := db.NewEvent("Coord__postStart_loadLiveStream")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{}
	stmt := `SELECT * FROM stream_item;`

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
			db.LogError(dbev, err, stmt, args)
			return nil, dbev
		}

		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	return entities, dbev
}
