package main

import (
	"fmt"

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
