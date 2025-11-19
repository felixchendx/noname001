package sqlite

func (db *DB) SysUser__SetPassword(userID, passwordHash string) (*DBEvent) {
	dbev := db.NewEvent("SysUser__SetPassword")

	args := []any{passwordHash, userID}
	stmt := `UPDATE sys_user SET password = ? WHERE id = ?;`

	if err := db.ExecuteWithArgs(stmt, args); err != nil {
		db.LogError(dbev, err, stmt, args)
		return dbev
	}

	return dbev
}
