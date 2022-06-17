package endpoints

import (
	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (e *Endpoints) PostAddDomainToQueue(ctx *fiber.Ctx) error {
	type schema struct {
		Domain string `json:"domain" validate:"required,domain,max=253"`
	}
	type response struct {
		ID *uuid.UUID `json:"id"`
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
	return ctx.JSON(&response{ID: id})
}
