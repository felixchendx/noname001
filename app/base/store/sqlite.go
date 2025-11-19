package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"noname001/logging"
)

const (
	SQLITE_QUERY_TIMEOUT = 10 * time.Second

	SQLITE_DATETIME_FORMAT = "2006-01-02 15:04:05"
)

var (
	// errors in this lib, see: https://pkg.go.dev/database/sql#pkg-variables
	SQLiteErrConnDone = errors.New("sqlite: connection is already closed")
	SQLiteErrNoRows   = errors.New("sqlite: no rows in result set")
	SQLiteErrTxDone   = errors.New("sqlite: transaction has already been committed or rolled back")
)

type SQLiteDBParams struct {
	Context     context.Context
	Logger      *logging.WrappedLogger
	LogPrefix   string

	SQLiteDBFileLocation string
}
type SQLiteDB struct {
	Context     context.Context
	Cancel      context.CancelFunc

	Logger      *logging.WrappedLogger
	LogPrefix   string

	SQLiteDBFileLocation string
	Conn                 *sql.DB
}

func NewSQLiteDB(params SQLiteDBParams) (*SQLiteDB, error) {
	sqliteDB := &SQLiteDB{}
	sqliteDB.Context, sqliteDB.Cancel = context.WithCancel(params.Context)
	sqliteDB.Logger = params.Logger
	sqliteDB.LogPrefix = params.LogPrefix

	sqliteDB.SQLiteDBFileLocation = params.SQLiteDBFileLocation
	dbConnString := fmt.Sprintf(
		"file://%s?_pragma=foreign_keys(0)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)",
		sqliteDB.SQLiteDBFileLocation,
	)

	dbConn, err := sql.Open("sqlite", dbConnString)
	if err != nil {
		return nil, err
	}

	sqliteDB.Conn = dbConn

	return sqliteDB, nil
}

func (db *SQLiteDB) Close() {
	db.Conn.Close()
	db.Cancel()
}

func (db *SQLiteDB) Execute(sqlStatement string) (error) {
	queryContext, queryCancel := context.WithTimeout(db.Context, SQLITE_QUERY_TIMEOUT)
	defer queryCancel()

	dbTransaction, err := db.Conn.BeginTx(queryContext, db.DefaultTransactionOptions())
	if err != nil {
		// If a non-default isolation level is used that the driver doesn't support,
		// an error will be returned.
		return err
	}
	defer dbTransaction.Rollback()

	result, err := dbTransaction.Exec(sqlStatement)
	if err != nil {
		return err
	}
	_ = result

	if err := dbTransaction.Commit(); err != nil {
		// ErrTxDone | ErrBadConn
		return err
	}

	return nil
}

// both positional(?) and named(@nwamed) can be used
// but the args must be typed NamedArgs or using the convenient func Named("nwamed", "vwalue")
// see: https://pkg.go.dev/database/sql#Named
func (db *SQLiteDB) ExecuteWithArgs(sqlStatement string, args []any) (error) {
	qctx, qc := context.WithTimeout(db.Context, SQLITE_QUERY_TIMEOUT)
	defer qc()

	tx, err := db.Conn.BeginTx(qctx, db.DefaultTransactionOptions())
	if err != nil {
		// If a non-default isolation level is used that the driver doesn't support,
		// an error will be returned.
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(sqlStatement, args...) 
	if err != nil {
		return err
	}
	_ = result

	if err := tx.Commit(); err != nil {
		// ErrTxDone | ErrBadConn
		return err
	}

	return nil
}

func (db *SQLiteDB) DefaultTransactionOptions() (*sql.TxOptions) {
	return &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
}

// ============================ event and logging =========================== //
func (db *SQLiteDB) NewEvent(name string) (*PersistenceEvent) {
	return NewPersistenceEvent(name)
}

func (db *SQLiteDB) LogError(pev *PersistenceEvent, err error, stmt string, args any) {
	pev.MarkAsError(err)

	db.Logger.Errorf(
		"%s::%s: __%s__ || %v || %s || %s",
		db.LogPrefix,
		pev.Name,
		pev.ID,
		pev.OriginalErr,
		stmt, args,
	)
	pev.Logged = true
}

func (db *SQLiteDB) LogWarning(pev *PersistenceEvent, format string, args ...any) {
	prefix := fmt.Sprintf("%s::%s: ", db.LogPrefix, pev.Name)
	
	db.Logger.Warnf(
		prefix + format,
		args...
	)
	pev.Logged = true
}

func (db *SQLiteDB) LogInfo(pev *PersistenceEvent, format string, args ...any) {
	prefix := fmt.Sprintf("%s::%s: ", db.LogPrefix, pev.Name)
	
	db.Logger.Infof(
		prefix + format,
		args...
	)
	pev.Logged = true
}

func (db *SQLiteDB) LogDebug(pev *PersistenceEvent, format string, args ...any) {
	prefix := fmt.Sprintf("%s::%s: ", db.LogPrefix, pev.Name)
	
	db.Logger.Debugf(
		prefix + format,
		args...
	)
	pev.Logged = true
}
// ============================ event and logging =========================== //

// ================================ converter =============================== //
func (db *SQLiteDB) BoolToString(b bool) string {
	if b { return "true" }
	return "false"
}

func (db *SQLiteDB) StringToBool(s string) bool {
	if s == "true" { return true }
	return false
}
// ================================ converter =============================== //