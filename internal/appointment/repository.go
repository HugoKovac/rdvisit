package appointment

import (
	"context"
	"time"

	"github.com/HugoKovac/rdvisit/db/sqlc"
	"github.com/HugoKovac/rdvisit/internal/domain"
	"github.com/HugoKovac/rdvisit/pkg/errors"
	"github.com/google/uuid"
)

type IRepository interface {
	CreateAppointment(ctx context.Context, date time.Time, practitionerID uuid.UUID, patientEmail string) (*domain.Appointment, error)
}

// ==================================
//	Params
// ==================================

type CreateParams struct {
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
}

//==================================
//			Implementation
//==================================

type sqlcRepository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) IRepository {
	return &sqlcRepository{q: q}
}

func (r *sqlcRepository) CreateAppointment(ctx context.Context, date time.Time, practitionerID uuid.UUID, patientEmail string) (*domain.Appointment, error) {
	appointment, err := r.q.CreateAppointment(ctx, sqlc.CreateAppointmentParams{
		Date:           date,
		PractitionerID: practitionerID,
		PatientEmail:   patientEmail,
	})
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &domain.Appointment{
		ID:             appointment.ID,
		Date:           appointment.Date,
		PractitionerID: appointment.PractitionerID,
		PatientEmail:   appointment.PatientEmail,
		CreatedAt:      appointment.CreatedAt,
		UpdatedAt:      appointment.UpdatedAt,
	}, nil
}
