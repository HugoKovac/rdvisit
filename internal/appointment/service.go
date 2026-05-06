package appointment

import (
	"context"
	"time"

	"github.com/HugoKovac/rdvisit/internal/domain"
	"github.com/HugoKovac/rdvisit/pkg/errors"
	"github.com/google/uuid"
)

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAppointment(ctx context.Context, date time.Time, practitionerID uuid.UUID, patientEmail string) (*domain.Appointment, error) {
	centers, err := s.repo.CreateAppointment(ctx, date, practitionerID, patientEmail)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return centers, nil
}

// ==================================
//	Params
// ==================================
