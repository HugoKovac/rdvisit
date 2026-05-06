package center

import (
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

func (h *Handler) Register(app *fiber.App, authMiddleware fiber.Handler, centerMiddleware fiber.Handler) {
	g := app.Group("/centers", authMiddleware)
	g.Get("/", h.GetCenters)
	g.Get("/:center_id", centerMiddleware, h.GetCenter)
}

//==================================
//				DTO
//==================================

type CentersResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Street       string    `json:"street"`
	StreetNumber string    `json:"street_number"`
	City         string    `json:"city"`
	PostalCode   string    `json:"postal_code"`
	Region       string    `json:"region"`
	Country      string    `json:"country"`
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

func (h *Handler) GetCenter(c fiber.Ctx) error {
	center, err := fibercontext.GetCenter(c)
	if err != nil {
		return err
	}
	return c.JSON(CentersResponse{
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
