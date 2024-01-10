package db

import (
	"database/sql"
)

type Queryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Prepare(query string) (*sql.Stmt, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

var (
	ErrNoRows = sql.ErrNoRows
)

type DBConn struct {
	DB *sql.DB
}

func NewSqlDB(sqlDB *sql.DB) DBConn {
	return DBConn{
		DB: sqlDB,
	}
}

func (db DBConn) Exec(query string, args ...any) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

func (db DBConn) Query(query string, args ...any) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

func (db DBConn) QueryRow(query string, args ...any) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

func (db DBConn) Prepare(query string) (*sql.Stmt, error) {
	return db.DB.Prepare(query)
}

func (db DBConn) Begin() (DBTx, error) {
	tx, err := db.DB.Begin()
	return DBTx{tx}, err
}
