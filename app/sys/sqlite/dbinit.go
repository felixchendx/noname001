package sqlite

func (db *DB) DBInit() (*DBEvent) {
	dbev := db.NewEvent("DBInit")

	sqlStatement := `
		CREATE TABLE IF NOT EXISTS ping (
			id TEXT NOT NULL PRIMARY KEY
		);

		INSERT INTO ping (id) SELECT 'pong'
		WHERE NOT EXISTS(SELECT 1 FROM ping WHERE id = 'pong');

		CREATE TABLE IF NOT EXISTS sys (
			id         TEXT NOT NULL PRIMARY KEY,

			code       TEXT NOT NULL DEFAULT '',
			name       TEXT NOT NULL DEFAULT '',
			type       TEXT NOT NULL DEFAULT '',	-- 'app' | 'mod'
			version    TEXT NOT NULL DEFAULT '',
			db_version BIGINT NOT NULL DEFAULT 0,

			created_ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx__sys_code ON sys(code);
		CREATE INDEX IF NOT EXISTS idx__sys_type ON sys(type);

		INSERT INTO sys (id, code, name, type, version, db_version)
			SELECT 'a5f5fa11-0947-4e4f-a176-31c45b28667b', 'generic001', 'Generic001', 'app', '0.0.0', 0
		WHERE NOT EXISTS(
			SELECT 1 FROM sys WHERE id = 'a5f5fa11-0947-4e4f-a176-31c45b28667b'
		);

		CREATE TABLE IF NOT EXISTS sys_user (
			id TEXT NOT NULL PRIMARY KEY,

			username       TEXT NOT NULL DEFAULT '',
			password       TEXT NOT NULL DEFAULT '',
			role_simple    TEXT NOT NULL DEFAULT '',		-- temp, predefined: 'superadmin', 'admin', 'operator', 'viewer'

			UNIQUE(username)
		);
		CREATE INDEX IF NOT EXISTS idx__sys_user__username ON sys_user(username);

		INSERT INTO sys_user (id, username, password, role_simple)
			SELECT 'a5f5fa11-0947-4e4f-a176-31c45b28667b', 'superadmin', '$2a$13$4lUGZ5kdErToL4LzsaXL5eo7m8rAISTbKU9ojbOCQajHQ8SODEzFu', 'superadmin'
		WHERE NOT EXISTS(
			SELECT 1 FROM sys_user WHERE id = 'a5f5fa11-0947-4e4f-a176-31c45b28667b'
		);
	`
	if err := db.Execute(sqlStatement); err != nil {
		db.LogError(dbev, err, "", "")
		return dbev
	}

	return dbev
}