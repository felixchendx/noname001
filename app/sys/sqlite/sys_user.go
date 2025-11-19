package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (db *DB) SysUser__Search(sc *SysUser__SearchCriteria) (*SysUser__SearchResult, *DBEvent) {
	dbev := db.NewEvent("SysUser__Search")

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
		if v := sc.UsernameLike; v != "" {
			sharedWhere = append(sharedWhere, "su.username LIKE @username")
			sharedArgs  = append(sharedArgs, sql.Named("username", fmt.Sprintf("%%%s%%", v)))
		}

		if v := sc.RoleSimple; v != nil && len(v) > 0 && v[0] != "" {
			whereInParts := make([]string, 0)

			for idx, item := range v {
				indexedArg := fmt.Sprintf("%s_%d", "role_simple", idx)
				sharedArgs = append(sharedArgs, sql.Named(indexedArg, item))

				whereInParts = append(whereInParts, fmt.Sprintf("@%s", indexedArg))
			}
			sharedWhere = append(sharedWhere, "su.role_simple IN (" + strings.Join(whereInParts, ", ") + ")")
		}

		sharedWhere = append(sharedWhere, "su.role_simple NOT IN ('superadmin')")
	}

	mainStmt := `
		SELECT *
		FROM sys_user AS su
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		ORDER BY
			su.username ASC
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
	
	entities := make([]*SysUserPE, 0)
	for rows.Next() {
		var entity *SysUserPE = &SysUserPE{}
		if err := rows.Scan(
			entity.fullScan()...
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
		SELECT COUNT(su.id)
		FROM sys_user AS su
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

	sr := &SysUser__SearchResult{
		Data: entities,
		Pagination: pagination,
	}

	return sr, dbev
}

func (db *DB) SysUser__AtomicInsert(entity *SysUserPE) (*SysUserPE, *DBEvent) {
	dbev := db.NewEvent("SysUser__AtomicInsert")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.ID,
		&entity.Username, &entity.Password,
		&entity.RoleSimple,
	}
	stmt := `
		INSERT INTO sys_user (
				id,
				username, password,
				role_simple
			) VALUES (
				?,
				?, ?,
				?
			)
		RETURNING *
		;
	`
	
	retEntity := &SysUserPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
	); err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	return retEntity, dbev
}

func (db *DB) SysUser__Get(id string) (*SysUserPE, *DBEvent) {
	dbev := db.NewEvent("SysUser__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{id}
	stmt := `SELECT * FROM sys_user WHERE id = ?;`

	retEntity := &SysUserPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
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

func (db *DB) SysUser__GetByUsername(username string) (*SysUserPE, *DBEvent) {
	dbev := db.NewEvent("SysUser__GetByUsername")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{username}
	stmt := `SELECT * FROM sys_user WHERE username = ?;`

	retEntity := &SysUserPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
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

func (db *DB) SysUser__AtomicUpdate(entity *SysUserPE) (*SysUserPE, *DBEvent) {
	dbev := db.NewEvent("SysUser__AtomicUpdate")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.Password,
		&entity.RoleSimple,
	}
	args = append(args, &entity.ID)
	stmt := `
		UPDATE sys_user SET
			password = ?,
			role_simple = ?
		WHERE id = ?
		RETURNING *
		;
	`

	retEntity := &SysUserPE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		retEntity.fullScan()...
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

func (db *DB) SysUser__AtomicDelete(id string) (*DBEvent) {
	dbev := db.NewEvent("SysUser__AtomicDelete")

	args := []any{id}
	stmt := `DELETE FROM sys_user WHERE id = ?;`

	err := db.ExecuteWithArgs(stmt, args)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return dbev
	}

	return dbev
}
