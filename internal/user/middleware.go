package user

import (
	"github.com/HugoKovac/rdvisit/pkg/errors"
	"github.com/HugoKovac/rdvisit/pkg/fiber/fibercontext"
	"github.com/gofiber/fiber/v3"
)

func CheckUserVerified(svc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		u, err := fibercontext.GetUserClaims(c)
		if err != nil {
			return errors.Wrap(errors.Forbidden)
		}

		user, err := svc.GetUserByID(ctx, u.ID)
		if err != nil {
			return err
		}

		if !user.Verified {
			return errors.Wrap(errors.Forbidden)
		}

		return c.Next()
	}
}
