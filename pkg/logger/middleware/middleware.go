package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	pkgerrors "github.com/HugoKovac/rdvisit/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func parseRequestFields(c fiber.Ctx) slog.Attr {
	return slog.Group(
		"request",
		"method", c.Method(),
		"path", c.Path(),
		"status", c.Response().StatusCode(),
	)
}

func composeErrorMessage(c fiber.Ctx, reqErr error) {
	msg := "info"
	attrs := []slog.Attr{parseRequestFields(c)}
	if reqErr != nil {
		var st pkgerrors.IErrorWrapper
		if errors.As(reqErr, &st) {
			attrs = append(attrs,
				slog.Group("trace",
					slog.String("caller", fmt.Sprintf("%+v", st.FormatTrace())),
				),
			)
		}
	}

	level := slog.LevelInfo
	status := c.Response().StatusCode()
	if status >= 500 {
		level = slog.LevelError
	} else if status >= 400 {
		level = slog.LevelWarn
	}

	slog.LogAttrs(c.RequestCtx(), level, msg, attrs...)
}

type InvalidArgument struct {
	Field  string `json:"field,omitempty"`
	Value  any    `json:"value,omitempty"`
	Tag    string `json:"tag,omitempty"`
	Param  string `json:"param,omitempty"`
	Actual string `json:"actual,omitempty"`
}

func determineError(reqErr error) (status int, msg any) {
	var (
		valErrors  validator.ValidationErrors
		collErrors []InvalidArgument
	)
	switch {
	case errors.As(reqErr, &valErrors):
		for _, f := range valErrors {
			collErrors = append(collErrors, InvalidArgument{Field: f.StructField(), Value: f.Value(), Tag: f.Tag(), Param: f.Param(), Actual: f.ActualTag()})
		}
		status, msg = http.StatusUnprocessableEntity, collErrors
	case errors.Is(reqErr, jwt.ErrTokenExpired):
		status, msg = http.StatusUnauthorized, "token expired"
	case errors.Is(reqErr, pkgerrors.Unauthorized):
		status, msg = http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	case errors.Is(reqErr, pkgerrors.BadRequest):
		status, msg = http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
	case errors.Is(reqErr, pkgerrors.UnprocessableEntity):
		status, msg = http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity)
	case errors.Is(reqErr, pkgerrors.NotFound):
		status, msg = http.StatusNotFound, http.StatusText(http.StatusNotFound)
	default:
		status, msg = http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}
	return
}

func Logger(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Pre Request
		reqErr := c.Next()
		// Post Request
		if reqErr != nil {
			status, msg := determineError(reqErr)
			c.Status(status).JSON(msg)
		}

		composeErrorMessage(c, reqErr)

		return nil
	}
}
