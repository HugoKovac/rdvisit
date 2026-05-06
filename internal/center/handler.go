package center

import (
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
	g := app.Group("/centers", authMiddleware)
	g.Get("/", h.GetCenters)
}

//==================================
//				DTO
//==================================

type CentersResponse struct {
	ID           uuid.UUID
	Name         string
	Street       string
	StreetNumber string
	City         string
	PostalCode   string
	Region       string
	Country      string
}

//==================================
//				Handler
//==================================

func (h *Handler) GetCenters(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	centers, err := h.svc.GetCenters(ctx)
	if err != nil {
		return err
	}

	rtn := make([]CentersResponse, 0, len(centers))
	for _, center := range centers {
		rtn = append(rtn, CentersResponse{
			ID:           center.ID,
			Name:         center.Name,
			Street:       center.Street,
			StreetNumber: center.StreetNumber,
			City:         center.City,
			PostalCode:   center.PostalCode,
			Region:       center.Region,
			Country:      center.Country,
		})
	}

	return c.JSON(rtn)
}
