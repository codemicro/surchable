package main

import (
	"github.com/codemicro/surchable/coordinator/urls"
	"github.com/codemicro/surchable/internal/config"
	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func run() error {
	_, err := db.New()
	if err != nil {
		return errors.WithStack(err)
	}

	app := setupApp()

	serveAddr := config.Coordinator.ServeHost + ":" + config.Coordinator.ServePort

	log.Info().Msgf("running coordinator on %s", serveAddr)

	if err := app.Listen(serveAddr); err != nil {
		return errors.Wrap(err, "fiber server run failed")
	}

	return nil
}

func main() {
	config.InitLogging()
	if err := run(); err != nil {
		log.Error().Stack().Err(err).Msg("failed to run coordinator")
	}
}

func setupApp() *fiber.App {
	app := fiber.New()

	app.Get(urls.OK, func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]string{"status": "ok"})
	})

	return app
}
