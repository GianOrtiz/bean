package db

import "database/sql"

type DBTx struct {
	DB *sql.Tx
}

func (tx DBTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.DB.Exec(query, args)
}

func (tx DBTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.DB.Query(query, args)
}

func (tx DBTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.DB.QueryRow(query, args)
}

func (tx DBTx) Prepare(query string) (*sql.Stmt, error) {
	return tx.DB.Prepare(query)
}

func (tx DBTx) Rollback() error {
	return tx.DB.Rollback()
}

func (tx DBTx) Commit() error {
	return tx.DB.Commit()
}
