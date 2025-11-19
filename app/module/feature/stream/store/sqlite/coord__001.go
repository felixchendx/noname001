package sqlite

import (
	"context"
)

func (db *DB) Coord__GetStreamItemIDListByStreamProfile(streamProfileID string) ([]string, *DBEvent) {
	dbev := db.NewEvent("Coord__GetStreamItemIDListByStreamProfile")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{streamProfileID}
	stmt := `
		SELECT
			si.id
		FROM stream_item AS si
			LEFT JOIN stream_group   AS sg ON sg.id = si.stream_group_id
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

func (db *DB) Coord__GetStreamItemIDListByStreamGroup(streamGroupID string) ([]string, *DBEvent) {
	dbev := db.NewEvent("Coord__GetStreamItemIDListByStreamGroup")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{streamGroupID}
	stmt := `
		SELECT
			si.id
		FROM stream_item AS si
			LEFT JOIN stream_group AS sg ON sg.id = si.stream_group_id
		WHERE 1=1
			AND sg.id = ?
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
