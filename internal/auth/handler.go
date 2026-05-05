package auth

import (
	"net/http"
	"project/clean/pkg/errors"
	"project/clean/pkg/fiber/fibercontext"

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
	g := app.Group("/auth")
	g.Post("/register", h.Create)
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
	g.Post("/logout", h.Logout)
	g.Post("/verify", authMiddleware, h.Verify)
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

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=12,max=50"`
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

type VerifyRequest struct {
	Code string `json:"code" validate:"len=6"`
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

	user, err := h.svc.Register(c, registerRequest.Email, registerRequest.FirstName, registerRequest.LastName, registerRequest.Password)

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
		return err
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

func (h *Handler) Verify(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	u, err := fibercontext.GetUserClaims(c)
	if err != nil {
		return err
	}

	verifyRequest := &VerifyRequest{}
	if err := c.Bind().Body(verifyRequest); err != nil {
		return errors.Wrap(err)
	}

	if err := h.validate.Struct(verifyRequest); err != nil {
		return errors.Wrap(err)
	}

	verified, err := h.svc.Verify(ctx, u.ID, verifyRequest.Code)
	if err != nil {
		return err
	}

	if !verified {
		return errors.Wrap(errors.Unauthorized)
	}

	return nil
}
