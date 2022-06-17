package util

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type RichError struct {
	Status int
	Reason string
	Detail any
}

func NewRichError(status int, reason string, detail any) error {
	return &RichError{
		Status: status,
		Reason: reason,
		Detail: detail,
	}
}

func NewRichErrorFromFiberError(err *fiber.Error, detail any) error {
	return NewRichError(err.Code, err.Message, detail)
}

func (r *RichError) Error() string {
	return fmt.Sprintf("handler error, %d: %s", r.Status, r.Reason)
}

func (r *RichError) AsJSON() ([]byte, error) {
	info := map[string]any{
		"status":  "error",
		"message": r.Reason,
	}

	if r.Detail != nil {
		info["detail"] = r.Detail
	}

	return json.Marshal(info)
}

func JSONErrorHandler(ctx *fiber.Ctx, err error) error {
	var re *RichError
	if e, ok := err.(*fiber.Error); ok {
		re = NewRichErrorFromFiberError(e, nil).(*RichError)
	} else if e, ok := err.(*RichError); ok {
		re = e
	} else {
		log.Error().Stack().Err(err).Str("location", "fiber error handler").Str("route", ctx.OriginalURL()).Send()
		re = NewRichErrorFromFiberError(fiber.ErrInternalServerError, nil).(*RichError)
	}
	jsonBytes, err := re.AsJSON()
	if err != nil {
		jsonBytes = []byte(`{"status":"error","message":"Internal Server Error","detail":"unable to produce detailed description"}`)
		log.Error().Err(err).Str("location", "fiber error handler").Msg("unable to produce error response")
	}
	ctx.Status(re.Status)
	ctx.Type("json")
	return ctx.Send(jsonBytes)
}

func DetailFromValidationErrors(errs validator.ValidationErrors) []map[string]string {
	var o []map[string]string

	for _, err := range errs {
		x := make(map[string]string)
		x["field"] = err.Field()
		x["tag"] = err.Tag()
		x["param"] = err.Param()
		o = append(o, x)
	}

	return o
}
