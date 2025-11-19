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
		CREATE TABLE stream_profile
		(
			id                       TEXT NOT NULL PRIMARY KEY,
			code                     TEXT NOT NULL DEFAULT '',
			name                     TEXT NOT NULL DEFAULT '',
			state                    TEXT NOT NULL DEFAULT 'active',	-- 'active' | 'readonly' | 'inactive'
			note                     TEXT NOT NULL DEFAULT '',

			target_video_fps         REAL NOT NULL DEFAULT 0,
			target_video_width       INT NOT NULL DEFAULT 0,
			target_video_height      INT NOT NULL DEFAULT 0,
			target_video_codec       TEXT NOT NULL DEFAULT 'h264',		-- 'h264' | 'h265'
			target_video_compression INT NOT NULL DEFAULT 0,			-- relative to hardware information
			target_video_bitrate     INT NOT NULL DEFAULT 0,			-- unit in bit/s

			target_audio_codec       TEXT NOT NULL DEFAULT 'aac',		-- 'opus' | 'aac'
			target_audio_compression INT NOT NULL DEFAULT 0,
			target_audio_bitrate     INT NOT NULL DEFAULT 0,			-- unit in bit/s

			show_timestamp           TEXT NOT NULL DEFAULT 'none',		-- 'none' | 'all'
			show_video_info          TEXT NOT NULL DEFAULT 'none',		-- 'none' | 'all'
			show_audio_info          TEXT NOT NULL DEFAULT 'none',		-- 'none' | 'all'
			show_site_info           TEXT NOT NULL DEFAULT 'none',		-- 'none' | 'all'

			created_ts               DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts               DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(code)
		);
		CREATE INDEX idx__stream_profile__state ON stream_profile(state);
		
		INSERT INTO stream_profile (
			id,
			code, name, state, note,
			target_video_codec, target_video_compression, target_video_bitrate,
			target_audio_codec, target_audio_compression, target_audio_bitrate
		) VALUES (
			'06b047f9-3a67-46ec-8b78-40ef6406aee3',
			'H264_DEFAULT', 'h264 - Default Compression', 'readonly', 'h264 with default compression',
			'h264', 60, '300000',
			'aac',  60, '128000'
		);
		INSERT INTO stream_profile (
			id,
			code, name, state, note,
			target_video_codec, target_video_compression, target_video_bitrate,
			target_audio_codec, target_audio_compression, target_audio_bitrate
		) VALUES (
			'16f64f39-d98b-4805-9311-9df670720465',
			'H265_DEFAULT', 'h265 - Default Compression', 'readonly', 'h265 with default compression',
			'h265', 60, '300000',
			'aac',  60, '128000'
		);
		INSERT INTO stream_profile (
			id,
			code, name, state, note,
			target_video_codec, target_video_compression, target_video_bitrate,
			target_audio_codec, target_audio_compression, target_audio_bitrate
		) VALUES (
			'22cae1db-c60b-45d4-b306-8ca53c392682',
			'H264_RAW', 'h264 - No Compression', 'readonly', 'h264 without compression',
			'h264', 0, '0',
			'aac',  0, '0'
		);
		INSERT INTO stream_profile (
			id,
			code, name, state, note,
			target_video_codec, target_video_compression, target_video_bitrate,
			target_audio_codec, target_audio_compression, target_audio_bitrate
		) VALUES (
			'3a8bdc43-932d-477e-bc5e-4000e899e081',
			'H265_RAW', 'h265 - No Compression', 'readonly', 'h265 without compression',
			'h265', 0, '0',
			'aac',  0, '0'
		);

		CREATE TABLE stream_group
		(
			id                TEXT NOT NULL PRIMARY KEY,
			code              TEXT NOT NULL DEFAULT '',
			name              TEXT NOT NULL DEFAULT '',
			state             TEXT NOT NULL DEFAULT 'active',		-- 'active' | 'readonly' | 'inactive'
			note              TEXT NOT NULL DEFAULT '',
			
			stream_profile_id TEXT NOT NULL DEFAULT '',

			created_ts        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(code)
		);
		CREATE INDEX idx__stream_group__state             ON stream_group(state);
		CREATE INDEX idx__stream_group__stream_profile_id ON stream_group(stream_profile_id);

		CREATE TABLE stream_item
		(
			id                TEXT NOT NULL PRIMARY KEY,
			stream_group_id   TEXT NOT NULL DEFAULT '',
			code              TEXT NOT NULL DEFAULT '',
			name              TEXT NOT NULL DEFAULT '',
			state             TEXT NOT NULL DEFAULT 'active',		-- 'active' | 'readonly' | 'inactive'
			note              TEXT NOT NULL DEFAULT '', 

			source_type       TEXT NOT NULL DEFAULT '',		-- 'mod_device' | 'external' | 'file' | 'embedded_file'

			-- source_type 'mod_device'
			device_code        TEXT NOT NULL DEFAULT '',
			device_channel_id  TEXT NOT NULL DEFAULT '',
			device_stream_type TEXT NOT NULL DEFAULT '',

			-- source_type 'external'
			external_url      TEXT NOT NULL DEFAULT '',
	
			-- source_type 'file'
			filepath          TEXT NOT NULL DEFAULT '',

			-- source_type 'embedded_file'
			embedded_filepath TEXT NOT NULL DEFAULT '',

			created_ts        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			unique(code)
		);
		CREATE INDEX idx__stream_item__stream_group_id ON stream_item(stream_group_id);
		CREATE INDEX idx__stream_item__state           ON stream_item(state);
		CREATE INDEX idx__stream_item__source_type     ON stream_item(source_type);
	`); err != nil {
		db.LogError(dbev, err, "", updateVersion)
		return dbev
	}

	// TODO: 
	// - do not use 'read_only' flag on field state, use dedicated field 'defined_by' to indicate data ownership
	// - and then remove / reconfigure queries for system-defined data vs user-defined data

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
