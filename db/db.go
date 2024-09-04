package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPGStorage(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
