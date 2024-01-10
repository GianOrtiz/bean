package db

import (
	"database/sql"

	"github.com/GianOrtiz/bean/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func GetDBConnection(config config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.DatabaseConnectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
