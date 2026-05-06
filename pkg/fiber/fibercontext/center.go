package fibercontext

import (
	"github.com/HugoKovac/rdvisit/internal/domain"
	"github.com/HugoKovac/rdvisit/pkg/errors"

	"github.com/gofiber/fiber/v3"
)

func GetCenter(c fiber.Ctx) (domain.Center, error) {
	center, exists := c.Locals("center").(domain.Center)
	if !exists {
		return center, errors.Wrap(errors.NotFound)
	}
	return center, nil
}
