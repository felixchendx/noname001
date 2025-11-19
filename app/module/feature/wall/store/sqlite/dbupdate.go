package sqlite

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	dbUpdateErr001 = errors.New("missing critical information to do dbupdate.")
)

func (db *DB) DBUpdate() (*DBEvent) {
	dbev := db.NewEvent("DBUpdate")

	dbinfoPE, _dbev := db.getDBInfo()
	if _dbev.IsError() {
		return _dbev
	}
	if dbinfoPE == nil {
		db.LogError(dbev, dbUpdateErr001, "dbUpdateErr001", "")
		return dbev
	}

	currDBVersion := dbinfoPE.DBVersion
	var updateVersion int64 = 0


	updateVersion = 1
	if err := db.doDBUpdate(currDBVersion, updateVersion, `
		CREATE TABLE wall_layout
			/*
				.____.____.____.____.    .________.________.
				|_1__|_2__|_3__|_4__|    | 1      | 2      |
				|__5_|____|____|____|    |________|________|
				|____|____|____|____|    | 3      | 4      |
				|____|____|____|__16|    |________|________|

				^^^ 16 item of 1x1        ^^^ 4 item of 2x2
			*/
		(
			id                  TEXT NOT NULL PRIMARY KEY,
			code                TEXT NOT NULL DEFAULT '',
			name                TEXT NOT NULL DEFAULT '',
			state               TEXT NOT NULL DEFAULT 'active',		-- 'active' | 'inactive'
			note                TEXT NOT NULL DEFAULT '',

			layout_formation    TEXT NOT NULL DEFAULT '',
			layout_item_count   INT  NOT NULL DEFAULT 0,

			defined_by          TEXT NOT NULL DEFAULT 'user',		-- 'system' | 'user'
			created_ts          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(code)
		);
		CREATE INDEX idx__wall_layout__state ON wall_layout(state);

		INSERT INTO wall_layout (
			id,
			code, name, state, note,
			layout_formation, layout_item_count,
			defined_by
		) VALUES (
			'06b047f9-3a67-46ec-8b78-40ef6406aee3',
			'DEFAULT_4', 'Default layout - 4 stream', 'active', '',
			'4:2x2', 4,
			'system'
		);
		INSERT INTO wall_layout (
			id,
			code, name, state, note,
			layout_formation, layout_item_count,
			defined_by
		) VALUES (
			'2d1769e5-6999-44c0-a37a-3a8fac90b6fb',
			'DEFAULT_12', 'Default layout - 12 stream', 'active', '',
			'12:1x1', 12,
			'system'
		);
		INSERT INTO wall_layout (
			id,
			code, name, state, note,
			layout_formation, layout_item_count,
			defined_by
		) VALUES (
			'16f64f39-d98b-4805-9311-9df670720465',
			'DEFAULT_16', 'Default layout - 16 stream', 'active', '',
			'16:1x1', 16,
			'system'
		);
		INSERT INTO wall_layout (
			id,
			code, name, state, note,
			layout_formation, layout_item_count,
			defined_by
		) VALUES (
			'22cae1db-c60b-45d4-b306-8ca53c392682',
			'DEFAULT_1B7S', 'Default layout - 1 big stream, 7 small stream', 'active', '',
			'1:3x3_7:1x1', 8,
			'system'
		);

		CREATE TABLE wall
		(
			id             TEXT NOT NULL PRIMARY KEY,
			code           TEXT NOT NULL DEFAULT '',
			name           TEXT NOT NULL DEFAULT '',
			state          TEXT NOT NULL DEFAULT 'active',	-- 'active' | 'inactive'
			note           TEXT NOT NULL DEFAULT '',

			wall_layout_id TEXT NOT NULL DEFAULT '',

			created_ts     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(code)
		);
		CREATE INDEX idx__wall__state          ON wall(state);
		CREATE INDEX idx__wall__wall_layout_id ON wall(wall_layout_id);

		CREATE TABLE wall_item
		(
			id                TEXT NOT NULL PRIMARY KEY,
			wall_id           TEXT NOT NULL DEFAULT '',
			idx               INT NOT NULL DEFAULT 0,

			source_node_id    TEXT NOT NULL DEFAULT '',
			stream_code       TEXT NOT NULL DEFAULT '',

			created_ts        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(wall_id, idx)
		);
		CREATE INDEX idx__wall_item__source_node_id ON wall_item(source_node_id);
		CREATE INDEX idx__wall_item__stream_code    ON wall_item(stream_code);

		CREATE TABLE wall_rotation
		(
			id            TEXT NOT NULL PRIMARY KEY,
			code          TEXT NOT NULL DEFAULT '',
			name          TEXT NOT NULL DEFAULT '',
			state         TEXT NOT NULL DEFAULT 'active',	-- 'active' | 'inactive'
			note          TEXT NOT NULL DEFAULT '',

			item_count    INT NOT NULL DEFAULT 0,

			created_ts    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(code)
		);
		CREATE INDEX idx__wall_rotation__state ON wall_rotation(state);

		CREATE TABLE wall_rotation_item
		(
			id                  TEXT NOT NULL PRIMARY KEY,
			wall_rotation_id    TEXT NOT NULL DEFAULT '',
			idx                 INT NOT NULL DEFAULT 0,

			wall_id             TEXT NOT NULL DEFAULT '',
			display_time        INT NOT NULL DEFAULT 0, -- in seconds

			created_ts          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(wall_rotation_id, idx)
		);
		CREATE INDEX idx__wall_rotation_item__wall_id ON wall_rotation_item(wall_id);
	`); err != nil {
		db.LogError(dbev, err, "", updateVersion)
		return dbev
	}

	// TODO: wall layout ordering

	return dbev
}

func (db *DB) doDBUpdate(currDBVersion int64, dbUpdateVersion int64, sqlStatement string) (error) {
	dbev := db.NewEvent("doDBUpdate")
	
	if currDBVersion >= dbUpdateVersion {
		db.LogDebug(dbev, "[%04v] skipping... (%v >= %v)", dbUpdateVersion, currDBVersion, dbUpdateVersion)
		return nil
	}
	db.LogInfo(dbev, "[%04v] updating...", dbUpdateVersion)

	sqlStatement += fmt.Sprintf(
		`
		UPDATE dbinfo SET db_version = %s WHERE code = '%s';
		`,
		strconv.FormatInt(dbUpdateVersion, 10),
		MODULE_CODE,
	)

	err := db.Execute(sqlStatement)

	return err
}
