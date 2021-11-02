package server

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type DBConfig struct {
	PGHost     string `name:"pghost" env:"PGHOST"`
	PGDatabase string `name:"pgdb" env:"PGDATABASE"`
	PGUser     string `name:"pguser" env:"PGUSER"`
	PGPassword string `name:"pgpassword" env:"PGPASSWORD"`
	PGSSLMode  string `name:"pgsslmode" env:"PGSSLMODE"`
	PGPort     int    `name:"pgport" env:"PGPORT" value:"5432"`
}

func (c *DBConfig) url() string {
	u := url.URL{
		User:   url.UserPassword(c.PGUser, c.PGPassword),
		Host:   c.PGHost + ":" + strconv.Itoa(c.PGPort),
		Scheme: "postgres",
		Path:   c.PGDatabase,
	}
	q := u.Query()
	q.Set("sslmode", c.PGSSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func initDB(conf DBConfig) (*Queries, error) {
	cfg, err := pgx.ParseConfig(conf.url())
	if err != nil {
		return nil, err
	}
	conn := stdlib.OpenDB(*cfg)
	err = validateSchema(conn)
	if err != nil {
		return nil, fmt.Errorf("validating schema: %w", err)
	}
	db := New(conn)
	return db, nil
}
