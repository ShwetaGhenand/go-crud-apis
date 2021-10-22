package server

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func initDB(url string) *sql.DB {
	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		log.Fatalf("Error connecting database: %v.", err)
	}
	db := stdlib.OpenDB(*cfg)
	content, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		log.Fatalf("Error reading sql file: %v.", err)
	}
	query := string(content)
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Error creating table: %v.", err)
	}
	return db
}
