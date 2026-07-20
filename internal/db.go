package internal

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

var DBConn *sql.DB

func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error
	DBConn, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}
	return DBConn, DBConn.Ping()
}
