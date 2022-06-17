package endpoints

import (
	"github.com/gofiber/fiber/v2"
)

func (e *Endpoints) GetStatus(ctx *fiber.Ctx) error {
	return ctx.JSON(
		map[string]string{
			"status": "ok",
		},
	)
}
