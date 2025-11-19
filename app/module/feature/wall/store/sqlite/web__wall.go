package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Web__WallListingSearchCriteria struct {
	WallCodeLike string
	WallNameLike string
	WallState    string

	// OrderBy

	Pagination *SearchPagination
}
type Web__WallListingSearchResult struct {
	Data []map[string]string

	Pagination *SearchPagination
}

// type Web__WallListingItem struct {

// }

func (db *DB) Web__WallListing(sc *Web__WallListingSearchCriteria) (*Web__WallListingSearchResult, *DBEvent) {
	dbev := db.NewEvent("Web__WallListing")

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
	sharedArgs := []any{
		sql.Named("perpage", pagination.PerPage),
		sql.Named("page_offset", (pagination.CurrPage-1)*pagination.PerPage),
	}

	if sc != nil {
		if v := sc.WallCodeLike; v != "" {
			sharedWhere = append(sharedWhere, "w.code LIKE @wall_code")
			sharedArgs = append(sharedArgs, sql.Named("wall_code", fmt.Sprintf("%%%s%%", v)))
		}

		if v := sc.WallNameLike; v != "" {
			sharedWhere = append(sharedWhere, "w.name LIKE @wall_name")
			sharedArgs = append(sharedArgs, sql.Named("wall_name", fmt.Sprintf("%%%s%%", v)))
		}

		if v := sc.WallState; v != "" {
			sharedWhere = append(sharedWhere, "w.state = @wall_state")
			sharedArgs = append(sharedArgs, sql.Named("wall_state", v))
		}
	}

	mainStmt := `
		SELECT
			w.id                 AS wall_id,
			w.code               AS wall_code,
			w.name               AS wall_name,
			w.state              AS wall_state,
			w.note               AS wall_note,

			wl.code              AS wall_layout_code,
			wl.layout_item_count AS wall_layout_item_count,

			COUNT(wi.wall_id)    AS wall_item_set_count

		FROM wall AS w
			LEFT JOIN wall_layout AS wl ON wl.id = w.wall_layout_id
			LEFT JOIN wall_item   AS wi ON wi.wall_id = w.id AND wi.stream_code != ''
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		GROUP BY
			w.id
		ORDER BY
			w.code ASC
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

	colNames, err := rows.Columns()
	if err != nil {
		db.LogError(dbev, err, "mainQ", "rows.columns")
		return nil, dbev
	}

	dataRows := make([]map[string]string, 0)

	for rows.Next() {
		colDat := make([]string, len(colNames))
		colDatPointers := make([]any, len(colNames))
		for i, _ := range colNames {
			colDatPointers[i] = &colDat[i]
		}

		if err := rows.Scan(
			colDatPointers...,
		); err != nil {
			db.LogError(dbev, err, "mainQ", "rows.scan")
			return nil, dbev
		}

		dataRow := make(map[string]string, len(colNames))
		for i, colName := range colNames {
			dataRow[colName] = colDat[i]
		}

		dataRows = append(dataRows, dataRow)
	}

	if err := rows.Err(); err != nil {
		db.LogError(dbev, err, "mainQ", "rows.err")
		return nil, dbev
	}
	// === ^^^ Main ^^^ === //

	// === VVV Pagination VVV === //
	totalStmt := `
		SELECT COUNT(w.id)
		FROM wall AS w
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

	sr := &Web__WallListingSearchResult{
		Data:       dataRows,
		Pagination: pagination,
	}

	return sr, dbev
}
