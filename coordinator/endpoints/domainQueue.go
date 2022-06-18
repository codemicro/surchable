package endpoints

import (
	"fmt"

	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (e *Endpoints) Post_AddDomainToQueue(ctx *fiber.Ctx) error {
	type schema struct {
		Domain string `json:"domain" validate:"required,domain,max=253"`
	}
	type response struct {
		ID uuid.UUID `json:"id"`
	}

	inputData := new(schema)
	if err := util.ParseAndValidateJSONBody(ctx, inputData); err != nil {
		return err
	}

	id, err := e.db.DomainQueueInsert(inputData.Domain)
	if err != nil {
		if errors.Is(err, db.ErrDomainAlreadyQueued) {
			return util.NewRichError(fiber.StatusConflict, "domain already queued", nil)
		}
		return errors.WithStack(err)
	}

	ctx.Status(fiber.StatusAccepted)
	return ctx.JSON(&response{ID: *id})
}

func (e *Endpoints) Get_CrawlerRequestJob(ctx *fiber.Ctx) error {
	type response struct {
		ID     uuid.UUID `json:"id"`
		Domain string    `json:"domain"`
		Start  string    `json:"start"`
	}

	crawlerID := ctx.Get(headerCrawlerID)
	if crawlerID == "" {
		return util.NewRichError(fiber.StatusBadRequest, fmt.Sprintf("%s header missing", headerCrawlerID), nil)
	}

	createdJob, err := e.db.RequestJob(crawlerID)
	if err != nil {
		if errors.Is(err, db.ErrWorkerIDInUse) {
			return util.NewRichError(fiber.StatusConflict, "crawler ID in use", nil)
		} else if errors.Is(err, db.ErrNoQueuedDomains) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}
		return errors.WithStack(err)
	}

	queueItem, err := e.db.DomainQueueFetch(createdJob.QueueItem)
	if err != nil {
		return errors.WithStack(err)
	}

	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(&response{
		ID:     createdJob.ID,
		Domain: queueItem.Domain,
		Start:  "/",
	})
}
