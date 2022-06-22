package util

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func init() {
	domainRe := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$`)

	if err := validate.RegisterValidation("domain", func(f validator.FieldLevel) bool {
		s := f.Field().String()
		return (strings.Count(s, ".") > 1 && domainRe.MatchString(s)) || s == ""
	}); err != nil {
		log.Fatal().Err(err).Msg("could not register domain validator")
	}
}

func ParseAndValidateJSONBody[T any](ctx *fiber.Ctx, output *T) error {
	if ctx.Get(fiber.HeaderContentType) != "application/json" {
		return NewRichError(fiber.StatusBadRequest, "ContentType must be `application/json`", nil)
	}

	body := ctx.Body()
	if len(body) == 0 {
		return NewRichError(fiber.StatusBadRequest, "empty JSON body", nil)
	}

	if err := json.Unmarshal(body, output); err != nil {
		return NewRichError(fiber.StatusBadRequest, "could not parse JSON body", err)
	}

	if err := validate.Struct(output); err != nil {
		errs := err.(validator.ValidationErrors)
		return NewRichError(fiber.StatusBadRequest, "failed validation", DetailFromValidationErrors(errs))
	}
	return nil
}
