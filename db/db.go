package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func (t *Database) Connect(username string, password string, host string, database string) error {
	var err error

	if t.db != nil {
		return fmt.Errorf("database already connected")
	}

	t.db, err = sql.Open("mysql", username+":"+password+"@tcp("+host+")/"+database)
	if err != nil {
		return err
	}

	return t.db.Ping()
}

func (t *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.db.Exec(query, args...)
}

func (t *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.db.QueryRow(query, args...)
}

func (t *Database) CloseDb() error {
	return t.db.Close()
}
