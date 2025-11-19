package sqlite

import (
	"context"

	baseStore "noname001/app/base/store"
)

func (db *DB) Coord__loadAutoStartDevice() ([]string, *baseStore.PersistenceEvent) {
	dbev := db.NewEvent("Coord__loadAutoStartDevice")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{}
	stmt := `SELECT id FROM device WHERE 1=1;`

	rows, err := db.Conn.QueryContext(qctx, stmt, args...)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	items := make([]string, 0) // rows have len ?
	for rows.Next() {
		var item string
		if err := rows.Scan(
			&item,
		); err != nil {
			db.LogError(dbev, err, stmt, args)
			return nil, dbev
		}
		items = append(items, item)
	}

	return items, dbev
}
