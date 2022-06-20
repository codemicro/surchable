package endpoints

import (
	"time"

	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	preflightLoad = "LOAD"
	preflightSkip = "SKIP"
)

func (e *Endpoints) Post_RequestPreflightCheck(ctx *fiber.Ctx) error {
	type schema struct {
		URL string `json:"url" validate:"required,url"`
	}
	type response struct {
		Permission string `json:"permission"`
	}

	inputData := new(schema)
	if err := util.ParseAndValidateJSONBody(ctx, inputData); err != nil {
		return err
	}

	pageLoad, err := e.db.QueryPageLoadsByURL(inputData.URL)
	if err != nil {
		if errors.Is(err, db.ErrNoMatchingPageLoad) {
			goto respondLoad
		}
		return err
	}

	if pageLoad.NotLoadBefore.After(time.Now()) {
		return ctx.JSON(&response{preflightSkip})
	}
	
respondLoad:
	return ctx.JSON(&response{preflightLoad})
}
