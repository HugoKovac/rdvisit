package center

import (
	"github.com/HugoKovac/rdvisit/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CenterURI struct {
	CenterID string `uri:"center_id" validator:"required,uuid"`
}

func CheckNGetCenter(svc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		uri := &CenterURI{}
		if err := c.Bind().URI(uri); err != nil {
			return errors.Wrap(err)
		}

		if err := validator.New().Struct(uri); err != nil {
			return errors.Wrap(err)
		}

		centerID := uuid.MustParse(uri.CenterID)
		center, err := svc.GetCenterByID(ctx, centerID)
		if err != nil {
			return err
		}

		c.Locals("center", *center)
		return c.Next()
	}
}
