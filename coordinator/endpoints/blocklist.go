package endpoints

import (
	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) Post_AddDomainToBlocklist(ctx *fiber.Ctx) error {
	type schema struct {
		Domain string `json:"domain" validate:"required,domain,max=253"`
		Reason string `json:"reason" validate:"required,max=75"`
	}

	inputData := new(schema)
	if err := util.ParseAndValidateJSONBody(ctx, inputData); err != nil {
		return err
	}

	err := e.db.AddDomainToBlocklist(inputData.Domain, inputData.Reason)
	if err != nil {
		if errors.Is(err, db.ErrDomainAlreadyInBlocklist) {
			goto ok
		}
		return errors.WithStack(err)
	}

ok:
	ctx.Status(fiber.StatusNoContent)
	return nil
}
