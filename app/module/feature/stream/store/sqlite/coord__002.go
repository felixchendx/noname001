package sqlite

import (
	"context"
)

func (db *DB) Coord__GetStreamItemIDListByDeviceCode(deviceCode string) ([]string, *DBEvent) {
	dbev := db.NewEvent("Coord__GetStreamItemIDListByDeviceCode")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{deviceCode}
	stmt := `
		SELECT
			si.id
		FROM stream_item AS si
		WHERE 1=1
			AND device_code = ?
		;
	`

	rows, err := db.Conn.QueryContext(qctx, stmt, args...)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	items := make([]string, 0)
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
