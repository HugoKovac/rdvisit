package auth

import (
	"project/clean/pkg/errors"

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
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
	g.Post("/logout", h.Logout)
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

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
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

	body := &RefreshRequest{}
	if err := c.Bind().Body(body); err != nil {
		return errors.Wrap(errors.Unauthorized)
	}

	claims, err := h.svc.ValidateToken(body.RefreshToken)
	if err != nil {
		return errors.Wrap(errors.Unauthorized)
	}

	if err := h.svc.DeleteRefreshToken(ctx, body.RefreshToken); err != nil {
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

func (h *Handler) Logout(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	body := &LogoutRequest{}
	if err := c.Bind().Body(body); err != nil {
		return err
	}
	if err := h.svc.DeleteRefreshToken(ctx, body.RefreshToken); err != nil {
		return err
	}
	return nil
}
