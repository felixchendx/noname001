package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (db *DB) StreamProfile__AtomicInsert(entity *StreamProfilePE) (*StreamProfilePE, *DBEvent) {
	dbev := db.NewEvent("StreamProfile__AtomicInsert")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.ID,
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.TargetVideoCodec, &entity.TargetVideoCompression, &entity.TargetVideoBitrate,

		&entity.TargetAudioCodec, &entity.TargetAudioCompression, &entity.TargetAudioBitrate,
	}
	stmt := `
		INSERT INTO stream_profile (
				id,
				code, name, state,
				note,

				target_video_codec, target_video_compression, target_video_bitrate,

				target_audio_codec, target_audio_compression, target_audio_bitrate
			) VALUES (
				?,
				?, ?, ?,
				?,

				?, ?, ?,

				?, ?, ?
			)
		RETURNING *
		;
	`

	retEntity := &StreamProfilePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamProfile__EntityFullScan(retEntity)...,
	); err != nil {
		db.LogError(dbev, err, stmt, args)
		return nil, dbev
	}

	return retEntity, dbev
}

func (db *DB) StreamProfile__Get(id string) (*StreamProfilePE, *DBEvent) {
	dbev := db.NewEvent("StreamProfile__Get")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{id}
	stmt := `SELECT * FROM stream_profile WHERE id = ?;`

	retEntity := &StreamProfilePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamProfile__EntityFullScan(retEntity)...,
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

func (db *DB) StreamProfile__GetByCode(code string) (*StreamProfilePE, *DBEvent) {
	dbev := db.NewEvent("StreamProfile__GetByCode")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{code}
	stmt := `SELECT * FROM stream_profile WHERE code = ?;`

	retEntity := &StreamProfilePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamProfile__EntityFullScan(retEntity)...,
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

func (db *DB) StreamProfile__AtomicUpdate(entity *StreamProfilePE) (*StreamProfilePE, *DBEvent) {
	dbev := db.NewEvent("StreamProfile__AtomicUpdate")

	qctx, qc := context.WithTimeout(db.Context, QUERY_TIMEOUT)
	defer qc()

	args := []any{
		&entity.Code, &entity.Name, &entity.State,
		&entity.Note,

		&entity.TargetVideoCodec, &entity.TargetVideoCompression, &entity.TargetVideoBitrate,

		&entity.TargetAudioCodec, &entity.TargetAudioCompression, &entity.TargetAudioBitrate,
	}
	args = append(args, &entity.ID)
	stmt := `
		UPDATE stream_profile SET
			code = ?, name = ?, state = ?,
			note = ?,

			target_video_codec = ?, target_video_compression = ?, target_video_bitrate = ?,

			target_audio_codec = ?, target_audio_compression = ?, target_audio_bitrate = ?,

			updated_ts = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING *
		;
	`

	retEntity := &StreamProfilePE{}
	if err := db.Conn.QueryRowContext(qctx, stmt, args...).Scan(
		db.StreamProfile__EntityFullScan(retEntity)...,
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

func (db *DB) StreamProfile__AtomicDelete(id string) *DBEvent {
	dbev := db.NewEvent("StreamProfile__AtomicDelete")

	args := []any{id}
	stmt := `DELETE FROM stream_profile WHERE id = ?;`

	err := db.ExecuteWithArgs(stmt, args)
	if err != nil {
		db.LogError(dbev, err, stmt, args)
		return dbev
	}

	return dbev
}

func (db *DB) StreamProfile__Search(sc *StreamProfile__SearchCriteria) (*StreamProfile__SearchResult, *DBEvent) {
	dbev := db.NewEvent("StreamProfile__Search")

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

	if sc != nil {
		if sc.CodeLike != "" {
			sharedWhere = append(sharedWhere, "sp.code LIKE @code")
			sharedArgs = append(sharedArgs, sql.Named("code", fmt.Sprintf("%%%s%%", sc.CodeLike)))
		}

		if sc.NameLike != "" {
			sharedWhere = append(sharedWhere, "sp.name LIKE @name")
			sharedArgs = append(sharedArgs, sql.Named("name", fmt.Sprintf("%%%s%%", sc.NameLike)))
		}

		if value := sc.State; value != nil && len(value) > 0 {
			if len(value) == 1 && value[0] != "" {
				sharedWhere = append(sharedWhere, "sp.state = @state")
				sharedArgs = append(sharedArgs, sql.Named("state", sc.State[0]))

			} else if len(value) > 1  {
				whereInParts := make([]string, 0, len(value))

				for idx, item := range value {
					if item == "" {
						continue
					}

					indexedArg := fmt.Sprintf("%s_%d", "state", idx)
					sharedArgs = append(sharedArgs, sql.Named(indexedArg, item))

					whereInParts = append(whereInParts, fmt.Sprintf("@%s", indexedArg))
				}
				sharedWhere = append(sharedWhere, "sp.state IN ("+strings.Join(whereInParts, ", ")+")")
			}
		}
	}

	mainStmt := `
		SELECT *
		FROM stream_profile as sp
		WHERE ` + strings.Join(sharedWhere, "\n			AND ") + `
		ORDER BY
			sp.code ASC
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

	entities := make([]*StreamProfilePE, 0)
	for rows.Next() {
		var entity *StreamProfilePE = &StreamProfilePE{}
		if err := rows.Scan(
			db.StreamProfile__EntityFullScan(entity)...
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
		SELECT COUNT(sp.id)
		FROM stream_profile as sp
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

	sr := &StreamProfile__SearchResult{
		Data:       entities,
		Pagination: pagination,
	}

	return sr, dbev
}
