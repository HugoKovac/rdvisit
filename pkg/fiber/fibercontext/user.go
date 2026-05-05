package fibercontext

import (
	"project/clean/internal/domain"
	"project/clean/pkg/errors"

	"github.com/gofiber/fiber/v3"
)

func GetUserClaims(c fiber.Ctx) (domain.User, error) {
	u, exists := c.Locals("userClaims").(domain.User)
	if !exists {
		return u, errors.Wrap(errors.NotFound)
	}
	return u, nil
}
