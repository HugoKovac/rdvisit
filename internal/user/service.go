package user

import (
	"context"
	"strings"
	"time"

	"github.com/HugoKovac/rdvisit/internal/domain"
	"github.com/HugoKovac/rdvisit/pkg/errors"

	"github.com/google/uuid"
)

type Service struct {
	repo   IRepository
	mailer IMailer
}

func NewService(repo IRepository, mailer IMailer) *Service {
	return &Service{repo: repo, mailer: mailer}
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return u, err
}

func (s *Service) Verify(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	vc, err := s.repo.GetVerificationCodeByUser(ctx, userID)
	if err != nil {
		return false, err
	}
	if strings.Compare(vc.Code, code) != 0 && vc.ExpiresAt.After(time.Now()) {
		return false, errors.Wrap(errors.Unauthorized)
	}
	if err := s.repo.VerifyUserByID(ctx, userID); err != nil {
		return false, err
	}
	return true, nil
}

// ==================================
//	Params
// ==================================
