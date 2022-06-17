package endpoints

import (
	"github.com/codemicro/surchable/coordinator/endpoints/urls"
	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
)

type Endpoints struct {
	db *db.DB
}

func New(dbi *db.DB) *Endpoints {
	return &Endpoints{
		db: dbi,
	}
}

func (e *Endpoints) SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: util.JSONErrorHandler,
	})

	app.Get(urls.OK, e.GetStatus)

	app.Post(urls.AddDomainToCrawlQueue, e.PostAddDomainToQueue)

	return app
}
