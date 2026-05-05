package auth

import (
	"project/clean/internal/domain"
	"project/clean/pkg/errors"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type AuthHeader struct {
	Token string `header:"Authorization"`
}

func AuthMiddleware(svc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := &AuthHeader{}
		if err := c.Bind().Header(header); err != nil || !strings.HasPrefix(header.Token, "Bearer ") {
			return errors.Wrap(errors.Unauthorized)
		}
		tokenString := strings.TrimPrefix(header.Token, "Bearer ")

		token, err := svc.ValidateToken(tokenString)
		if err != nil {
			return errors.Wrap(err)
		}

		c.Locals("userClaims", domain.User{
			ID:   token.ID,
			Role: token.Role,
		})

		return c.Next()
	}
}
