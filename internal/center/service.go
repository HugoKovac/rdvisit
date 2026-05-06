package center

import (
	"context"

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

func (s *Service) GetCenters(ctx context.Context) ([]*domain.Center, error) {
	centers, err := s.repo.GetCenters(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return centers, nil
}

func (s *Service) GetCenterByID(ctx context.Context, id uuid.UUID) (*domain.Center, error) {
	center, err := s.repo.GetCenterByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return center, nil
}

// ==================================
//	Params
// ==================================
