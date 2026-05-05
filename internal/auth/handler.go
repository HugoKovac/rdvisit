package auth

import (
	"project/clean/pkg/errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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
	g := app.Group("/auth")
	g.Post("/", h.Login)
	g.Post("/refresh", h.Refresh)
}

//==================================
//				DTO
//==================================

type AuthHeader struct {
	Token string `header:"Authorization"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=50"`
}

//==================================
//				Handler
//==================================

func (h *Handler) Login(c fiber.Ctx) error {
	loginRequest := &LoginRequest{}
	if err := c.Bind().Body(loginRequest); err != nil {
		return errors.Wrap(err)
	}

	if err := h.validate.Struct(loginRequest); err != nil {
		return errors.Wrap(err)
	}

	accessToken, refreshToken, err := h.svc.Login(c.RequestCtx(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Refresh(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	header := &AuthHeader{}
	if err := c.Bind().Header(header); err != nil {
		return errors.Wrap(errors.Unauthorized)
	}

	token := strings.Split(header.Token, "Bearer ")
	if len(token) != 2 {
		return errors.Wrap(errors.Unauthorized)
	}

	claims, err := h.svc.ValidateRefreshToken(token[1])
	if err != nil {
		return errors.Wrap(errors.Unauthorized)
	}

	if err := h.svc.RevokeRefreshToken(ctx, token[1]); err != nil {
		return err
	}

	accessToken, refreshToken, err := h.svc.GenerateTokenPair(ctx, claims.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
