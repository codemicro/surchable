package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/codemicro/surchable/internal/config"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type DB struct {
	pool *sql.DB
}

func New() (*DB, error) {
	dsn := fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' sslmode='disable'", config.DB.User, config.DB.Password, config.DB.DatabaseName, config.DB.Host)
	log.Debug().Str("dsn", dsn).Msg("Connecting to PostgreSQL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "could not open SQL connection")
	}

	rtn := &DB{
		pool: db,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := rtn.pool.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "could not ping database")
	}

	return rtn, nil
}
