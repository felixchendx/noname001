package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (db *DB) WallLayout__Search(sc *WallLayout__SearchCriteria) (*WallLayout__SearchResult, *DBEvent) {
	dbev := db.NewEvent("WallLayout__Search")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	// === VVV Normalize VVV === //
	var pagination *SearchPagination

	if sc != nil && sc.Pagination != nil {
		pagination = &(*sc.Pagination)
	} else {
		pagination = &SearchPagination{}
	}

	pagination.Normalize()
	// === ^^^ Normalize ^^^ === //

	sharedWhere := []string{"1=1"}
	sharedArgs  := []any{
		sql.Named("perpage", pagination.PerPage),
		sql.Named("page_offset", (pagination.CurrPage - 1) * pagination.PerPage),
	}

	if sc != nil {
		if v := sc.CodeLike; v != nil && *v != "" {
			sharedWhere = append(sharedWhere, "wl.code LIKE @code")
			sharedArgs  = append(sharedArgs, sql.Named("code", fmt.Sprintf("%%%s%%", *v)))
		}

		if v := sc.NameLike; v != nil && *v != "" {
			sharedWhere = append(sharedWhere, "wl.name LIKE @name")
			sharedArgs  = append(sharedArgs, sql.Named("name", fmt.Sprintf("%%%s%%", *v)))
		}

		if v := sc.State; v != nil && len(v) > 0 {
			whereInParts := make([]string, 0, len(v))

			for idx, item := range v {
				indexedArg := fmt.Sprintf("%s_%d", "state", idx)
				sharedArgs = append(sharedArgs, sql.Named(indexedArg, item))

				whereInParts = append(whereInParts, fmt.Sprintf("@%s", indexedArg))
			}
			sharedWhere = append(sharedWhere, "wl.state IN (" + strings.Join(whereInParts, ", ") + ")")

			// visz:
			// sharedWhere = append(sharedWhere, "wl.state IN (@state_0, @state_1)")
			// sharedArgs  = append(sharedArgs, sql.Named("state_0", "active"))
			// sharedArgs  = append(sharedArgs, sql.Named("state_1", "deprecated"))
		}
	}

	mainStmt := `
		SELECT *
		FROM wall_layout AS wl
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		ORDER BY
			wl.code ASC
		LIMIT @perpage
		OFFSET @page_offset
		;
	`

	// === VVV Main VVV === //
	rows, err := db.Conn.QueryContext(qctx, mainStmt, sharedArgs...)
	if err != nil {
		db.LogError(dbev, err, mainStmt, sharedArgs)
		return nil, dbev
	}
	
	entities := make([]*WallLayoutPE, 0)
	for rows.Next() {
		var entity *WallLayoutPE = &WallLayoutPE{}
		if err := rows.Scan(
			db.WallLayout__EntityFullScan(entity)...
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
	// === ^^^ Main ^^^ === //

	// === VVV Pagination VVV === //
	totalStmt := `
		SELECT COUNT(wl.id)
		FROM wall_layout AS wl
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		;
	`
	if err := db.Conn.QueryRowContext(qctx, totalStmt, sharedArgs...).Scan(
		&pagination.TotalData,
	); err != nil {
		db.LogError(dbev, err, totalStmt, sharedArgs)
		return nil, dbev
	}

	pagination.Finalize()
	// === ^^^ Pagination ^^^ === //

	sr := &WallLayout__SearchResult{
		Data: entities,
		Pagination: pagination,
	}

	return sr, dbev
}

func (db *DB) WallLayout__Get(id string) (*WallLayoutPE, *DBEvent) {
	dbev := db.NewEvent("WallLayout__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{id}
	stmt := `SELECT * FROM wall_layout WHERE id = ?;`

	retEntity := &WallLayoutPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.WallLayout__EntityFullScan(retEntity)...
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

func (db *DB) WallLayout__GetByCode(code string) (*WallLayoutPE, *DBEvent) {
	dbev := db.NewEvent("WallLayout__GetByCode")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{code}
	stmt := `SELECT * FROM wall_layout WHERE code = ?;`

	retEntity := &WallLayoutPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.WallLayout__EntityFullScan(retEntity)...
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
