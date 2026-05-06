package center

import (
	"context"

	"github.com/HugoKovac/rdvisit/internal/domain"
	"github.com/HugoKovac/rdvisit/pkg/errors"
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
	return centers, err
}

// ==================================
//	Params
// ==================================
