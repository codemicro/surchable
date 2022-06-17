package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/codemicro/surchable/internal/config"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type DB struct {
	pool *sql.DB
}

const maxConnectionAttempts = 4

func New() (*DB, error) {
	dsn := fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' sslmode='disable'", config.DB.User, config.DB.Password, config.DB.DatabaseName, config.DB.Host)
	log.Info().Msg("connecting to postgresql")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "could not open SQL connection")
	}

	rtn := &DB{
		pool: db,
	}

	for i := 1; i <= maxConnectionAttempts; i += 1 {
		logger := log.With().Int("attempt", i).Int("maxAttempts", maxConnectionAttempts).Logger()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := rtn.pool.PingContext(ctx)

		if err == nil {
			cancel()
			break
		}

		if e, ok := err.(*net.OpError); ((ok && e.Op == "dial") || errors.Is(err, context.DeadlineExceeded)) && i != maxConnectionAttempts {
			cancel()

			retryIn := int(math.Pow(math.E, float64(i)))
			logger.Warn().Err(err).Msgf("could not connect to database - retrying in %d seconds", retryIn)
			time.Sleep(time.Second * time.Duration(retryIn))

			continue
		}

		cancel()
		return nil, errors.Wrapf(err, "could not ping database after %d attempts", i)
	}

	return rtn, nil
}

func smartRollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Warn().Stack().Err(errors.WithStack(err)).Str("location", "smartRollback").Msg("failed to rollback transaction")
	}
}
