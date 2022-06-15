package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codemicro/surchable/internal/config"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type DB struct {
	db *sql.DB
}

func New() (*DB, error) {
	dsn := fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' sslmode='disable'", config.DB.User, config.DB.Password, config.DB.DatabaseName, config.DB.Host)
	log.Debug().Str("dsn", dsn).Msg("Connecting to PostgreSQL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func (db *DB) MakeConn() (*sql.Conn, error) {
	return db.db.Conn(context.Background())
}