package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Web__StreamGroupListingSearchCriteria struct {
	CodeLike          string
	NameLike          string
	State             []string

	StreamProfileLike string

	StreamItemLike    []string

	Pagination        *SearchPagination
}

type Web__StreamGroupListingSearchResult struct {
	Data       []map[string]string

	Pagination *SearchPagination
}

func (db *DB) Web__StreamGroupListing(sc *Web__StreamGroupListingSearchCriteria) (*Web__StreamGroupListingSearchResult, *DBEvent) {
	dbev := db.NewEvent("Web__StreamGroupListing")

	qCtx, qCancel := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qCancel()

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
		sql.Named("page_offset", (pagination.CurrPage - 1) * pagination.PerPage),
	}

	if sc != nil {
		if sc.CodeLike != "" {
			sharedWhere = append(sharedWhere, "sg.code LIKE @sg_code")
			sharedArgs = append(sharedArgs, sql.Named("sg_code", fmt.Sprintf("%%%s%%", sc.CodeLike)))
		}

		if sc.NameLike != "" {
			sharedWhere = append(sharedWhere, "sg.name LIKE @sg_name")
			sharedArgs = append(sharedArgs, sql.Named("sg_name", fmt.Sprintf("%%%s%%", sc.NameLike)))
		}

		if value := sc.State; value != nil && len(value) > 0 {
			whereInParts := make([]string, 0, len(value))

			for idx, item := range value {
				if item == "" {
					continue
				}

				indexedArg := fmt.Sprintf("%s_%d", "state", idx)
				sharedArgs = append(sharedArgs, sql.Named(indexedArg, item))

				whereInParts = append(whereInParts, fmt.Sprintf("@%s", indexedArg))
			}

			if len(whereInParts) > 0 {
				sharedWhere = append(sharedWhere, "sg.state IN ("+strings.Join(whereInParts, ", ")+")")
			}
		}
	}

	mainStmt := `
		SELECT
			sg.id        AS sg_id,
			sg.code      AS sg_code,
			sg.name      AS sg_name,
			sg.state     AS sg_state,
			sg.note      AS sg_note,

			sp.code      AS sp_code,

			count(si.id) AS si_count

		FROM stream_group AS sg
			LEFT JOIN stream_profile AS sp ON sp.id = sg.stream_profile_id
			LEFT JOIN stream_item    AS si ON si.stream_group_id = sg.id
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		GROUP BY
			sg.id
		ORDER BY
			sg.code ASC
		LIMIT @perpage
		OFFSET @page_offset
		;
	`

	// === VVV Main VVV === //
	rows, err := db.Conn.QueryContext(qCtx, mainStmt, sharedArgs...)
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
		SELECT COUNT(sg.id)
		FROM stream_group as sg
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		;
	`
	if err := db.Conn.QueryRowContext(qCtx, totalStmt, sharedArgs...).Scan(
		&pagination.TotalData,
	); err != nil {
		db.LogError(dbev, err, totalStmt, sharedArgs)
		return nil, dbev
	}

	pagination.Finalize()
	// === ^^^ Pagination ^^^ === //

	sr := &Web__StreamGroupListingSearchResult{
		Data:       dataRows,
		Pagination: pagination,
	}

	return sr, dbev
}
