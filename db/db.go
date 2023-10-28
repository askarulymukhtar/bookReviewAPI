package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "bookReview"
)

var DB *sql.DB

func OpenConnection() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	return nil
}

func CloseConnection() error {
	return DB.Close()
}
