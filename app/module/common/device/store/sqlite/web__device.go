package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Web__DeviceListingSearchCriteria struct {
	DeviceCodeLike     string
	DeviceNameLike     string
	DeviceState        string
	
	DeviceHostnameLike string
	DeviceUsernameLike string
	DeviceBrand        string
	Pagination         *SearchPagination
}

type Web__DeviceListingSearchResult struct {
	Data       []map[string]string

	Pagination *SearchPagination
}

func (db *DB) Web__DeviceListing(sc *Web__DeviceListingSearchCriteria) (*Web__DeviceListingSearchResult, *DBEvent) {
	dbev := db.NewEvent("Web__DeviceListing")

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
		sql.Named("page_offset", (pagination.CurrPage-1)*pagination.PerPage),
	}

	if (sc != nil) && (sc.DeviceCodeLike != "") {
		sharedWhere = append(sharedWhere, "code LIKE @device_code")
		sharedArgs = append(sharedArgs, sql.Named("device_code", fmt.Sprintf("%%%s%%", sc.DeviceCodeLike)))
	}

	if (sc != nil) && (sc.DeviceNameLike != "") {
		sharedWhere = append(sharedWhere, "name LIKE @device_name")
		sharedArgs = append(sharedArgs, sql.Named("device_name", fmt.Sprintf("%%%s%%", sc.DeviceNameLike)))
	}

	if (sc != nil) && (sc.DeviceState != "") {
		sharedWhere = append(sharedWhere, "state = @device_state")
		sharedArgs = append(sharedArgs, sql.Named("device_state", sc.DeviceState))
	}

	if (sc != nil) && (sc.DeviceHostnameLike != "") {
		sharedWhere = append(sharedWhere, "hostname LIKE @device_hostname")
		sharedArgs = append(sharedArgs, sql.Named("device_hostname", fmt.Sprintf("%%%s%%", sc.DeviceHostnameLike)))
	}

	if (sc != nil) && (sc.DeviceUsernameLike != "") {
		sharedWhere = append(sharedWhere, "username LIKE @device_username")
		sharedArgs = append(sharedArgs, sql.Named("device_username", fmt.Sprintf("%%%s%%", sc.DeviceUsernameLike)))
	}

	if (sc != nil) && (sc.DeviceBrand != "") {
		sharedWhere = append(sharedWhere, "brand = @device_brand")
		sharedArgs = append(sharedArgs, sql.Named("device_brand", sc.DeviceBrand))
	}

	mainStmt := `
		SELECT 
			id       AS device_id,
			code     AS device_code,
			name     AS device_name,
			state    AS device_state,
			note     AS device_note,

			protocol AS device_protocol,
			hostname AS device_hostname,
			port     AS device_port,
			username AS device_username,
			brand    AS device_brand
		FROM device
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		ORDER BY
			code ASC
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
		SELECT COUNT(id)
		FROM device
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

	sr := &Web__DeviceListingSearchResult{
		Data:       dataRows,
		Pagination: pagination,
	}

	return sr, dbev
}
