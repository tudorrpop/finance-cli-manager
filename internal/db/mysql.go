package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	dsn := "root:rootpassword@tcp(127.0.0.1:3306)/climanager"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
