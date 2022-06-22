package main

import (
	"fmt"
	"time"

	"github.com/codemicro/surchable/coordinator/endpoints"
	"github.com/codemicro/surchable/internal/config"
	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func run() error {
	database, err := db.New()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := database.Migrate(); err != nil {
		return errors.Wrap(err, "failed migration")
	}

	startTimeoutWorker(database, time.Minute*10)

	e := endpoints.New(database)
	app := e.SetupApp()

	serveAddr := config.Coordinator.ServeHost + ":" + config.Coordinator.ServePort

	log.Info().Msgf("starting coordinator server on %s", serveAddr)

	if err := app.Listen(serveAddr); err != nil {
		return errors.Wrap(err, "fiber server run failed")
	}

	return nil
}

func main() {
	config.InitLogging()
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		log.Error().Stack().Err(err).Msg("failed to run coordinator")
	}
}

func startTimeoutWorker(database *db.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			log.Info().Msg("running job timeout worker")
			err := database.RemoveTimedOutJobs()
			if err != nil {
				log.Error().Err(err).Str("location", "timeoutWorker").Send()
			}
		}
	}()
}
