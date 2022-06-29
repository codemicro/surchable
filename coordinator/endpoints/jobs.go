package endpoints

import (
	"fmt"
	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) Post_CompleteJob(ctx *fiber.Ctx) error {
	crawlerID := ctx.Get(headerCrawlerID)
	if crawlerID == "" {
		return util.NewRichError(fiber.StatusBadRequest, fmt.Sprintf("%s header missing", headerCrawlerID), nil)
	}

	err := e.db.CompleteJobByCrawlerID(crawlerID)
	if err != nil {
		if errors.Is(err, db.ErrNoActiveJob) {
			ctx.Status(fiber.StatusConflict)
			return nil
		}
		return errors.WithStack(err)
	}

	ctx.Status(fiber.StatusOK)
	return nil
}

func (e *Endpoints) Post_CancelJob(ctx *fiber.Ctx) error {
	crawlerID := ctx.Get(headerCrawlerID)
	if crawlerID == "" {
		return util.NewRichError(fiber.StatusBadRequest, fmt.Sprintf("%s header missing", headerCrawlerID), nil)
	}

	err := e.db.CancelJobByCrawlerID(crawlerID)
	if err != nil {
		if errors.Is(err, db.ErrNoActiveJob) {
			ctx.Status(fiber.StatusConflict)
			return nil
		}
		return errors.WithStack(err)
	}

	ctx.Status(fiber.StatusOK)
	return nil
}
