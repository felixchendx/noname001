package sqlite

import (
	"context"
)

func (db *DB) StreamGroup__GetByStreamProfileRef(streamProfileID string) ([]*StreamGroupPE, *DBEvent) {
	dbev := db.NewEvent("StreamGroup__GetByStreamProfileRef")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{streamProfileID}
	stmt := `
		SELECT
			sg.*
		FROM stream_group AS sg
			LEFT JOIN stream_profile AS sp ON sp.id = sg.stream_profile_id
		WHERE 1=1
			AND sp.id = ?
		;
	`

	rows, err := db.Conn.QueryContext(qctx, stmt, args...)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	items := make([]*StreamGroupPE, 0)
	for rows.Next() {
		var item *StreamGroupPE = &StreamGroupPE{}
		if err := rows.Scan(
			item.fullScan()...
		); err != nil {
			db.LogError(dbev, err, "item", "rows.scan")
			return nil, dbev
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		db.LogError(dbev, err, "item", "rows.err")
		return nil, dbev
	}

	return items, dbev
}
