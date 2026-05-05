package user

import (
	"context"
	"project/clean/internal/domain"
	"project/clean/pkg/errors"

	"github.com/google/uuid"
)

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return u, err
}

// ==================================
//	Params
// ==================================
