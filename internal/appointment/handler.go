package appointment

import (
	"time"

	"github.com/HugoKovac/rdvisit/pkg/errors"
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

func (h *Handler) Register(app *fiber.App, authMiddleware fiber.Handler, verifiedMiddleware fiber.Handler) {
	g := app.Group("/appointments", authMiddleware)
	g.Post("/", verifiedMiddleware, h.CreateAppointment)
}

//==================================
//				DTO
//==================================

type CreateAppointmentRequest struct {
	Date           string    `json:"date" validate:"datetime=2006-01-02T15:04:05Z02:00"`
	PractitionerID uuid.UUID `json:"practitioner_id"`
	PatientEmail   string    `json:"patient_email" validate:"email"`
}

type CreateAppointmentResponse struct {
	ID             uuid.UUID `json:"id"`
	Date           time.Time `json:"date"`
	PractitionerID uuid.UUID `json:"practitioner_id"`
	PatientEmail   string    `json:"patient_email" `
}

//==================================
//				Handler
//==================================

func (h *Handler) CreateAppointment(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	body := &CreateAppointmentRequest{}
	if err := c.Bind().Body(body); err != nil {
		return errors.Wrap(err)
	}

	if err := h.validate.Struct(body); err != nil {
		return errors.Wrap(err)
	}

	parsedDate, err := time.Parse("2006-01-02T15:04:05Z02:00", body.Date)
	if err != nil {
		return errors.Wrap(err)
	}

	appointment, err := h.svc.CreateAppointment(ctx, parsedDate, body.PractitionerID, body.PatientEmail)
	if err != nil {
		return err
	}
	return c.JSON(CreateAppointmentResponse{
		ID:             appointment.ID,
		Date:           appointment.Date,
		PractitionerID: appointment.PractitionerID,
		PatientEmail:   appointment.PatientEmail,
	})
}
