package endpoints

import (
	"fmt"
	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/urls"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
)

const headerCrawlerID = "X-Crawler-ID"

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

	app.Post(urls.AddDomainToCrawlQueue, e.Post_AddDomainToQueue)
	app.Get(urls.CrawlerRequestJob, e.Get_CrawlerRequestJob)
	app.Post(urls.CancelJob, e.Post_CancelJob)
	app.Post(urls.CompleteJob, e.Post_CompleteJob)

	app.Post(urls.RequestPreflightCheck, e.Post_RequestPreflightCheck)
	app.Post(urls.DigestPageLoad, e.Post_DigestPageLoad)

	app.Post(urls.BlocklistAdd, e.Post_AddDomainToBlocklist)

	return app
}

func getCrawlerIDHeader(ctx *fiber.Ctx) (string, error) {
	crawlerID := ctx.Get(headerCrawlerID)
	if crawlerID == "" {
		return "", util.NewRichError(
			fiber.StatusBadRequest,
			fmt.Sprintf("%s header missing", headerCrawlerID),
			nil,
		)
	}
	return crawlerID, nil
}
