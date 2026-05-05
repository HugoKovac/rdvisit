package user

import (
	"net/http"
	"project/clean/pkg/errors"

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

func (h *Handler) Register(app *fiber.App) {
	g := app.Group("/users")
	g.Post("/", h.Create)
}

//==================================
//				DTO
//==================================

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=12,max=50"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
}

//==================================
//				Handler
//==================================

func (h *Handler) Create(c fiber.Ctx) error {
	registerRequest := &RegisterRequest{}
	if err := c.Bind().Body(registerRequest); err != nil {
		return errors.Wrap(err)
	}

	if err := h.validate.Struct(registerRequest); err != nil {
		return errors.Wrap(err)
	}

	user, err := h.svc.Register(c, RegisterParams{
		registerRequest.Email,
		registerRequest.FirstName,
		registerRequest.LastName,
		registerRequest.Password,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(UserResponse{
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
	})
}
