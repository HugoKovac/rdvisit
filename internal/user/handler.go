package user

import (
	"github.com/HugoKovac/rdvisit/pkg/errors"
	"github.com/HugoKovac/rdvisit/pkg/fiber/fibercontext"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	svc      *Service
	validate *validator.Validate
}

func NewHandler(svc *Service) *Handler {
	validate := validator.New()
	return &Handler{svc: svc, validate: validate}
}

func (h *Handler) Register(app *fiber.App, authMiddleware fiber.Handler) {
	g := app.Group("/users", authMiddleware)
	g.Get("/me", h.Me)
}

//==================================
//				DTO
//==================================

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
}

//==================================
//				Handler
//==================================

func (h *Handler) Me(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	u, err := fibercontext.GetUserClaims(c)
	if err != nil {
		return err
	}

	fullUser, err := h.svc.GetUserByID(ctx, u.ID)
	if err != nil {
		errors.Wrap(errors.NotFound)
	}

	return c.JSON(UserResponse{
		fullUser.ID,
		fullUser.Email,
		fullUser.FirstName,
		fullUser.LastName,
	})
}
