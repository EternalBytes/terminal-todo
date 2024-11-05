package service

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./db/sqlite.db")
}
