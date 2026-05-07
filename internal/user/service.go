package user

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/HugoKovac/rdvisit/internal/domain"
	"github.com/HugoKovac/rdvisit/pkg/errors"

	"github.com/google/uuid"
)

type Service struct {
	repo                IRepository
	mailer              IMailer
	verificationCodeTTL time.Duration
}

func NewService(repo IRepository, mailer IMailer, verificationCodeTTL time.Duration) *Service {
	return &Service{
		repo:                repo,
		mailer:              mailer,
		verificationCodeTTL: verificationCodeTTL,
	}
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

func generateCode() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", errors.Wrap(err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *Service) CreateVerificationCode(ctx context.Context, userID uuid.UUID) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}

	if err := s.repo.CreateVerificationCode(ctx, userID, code, time.Now().Add(s.verificationCodeTTL)); err != nil {
		return "", err
	}

	return code, nil
}

func (s *Service) SendVerificationCode(ctx context.Context, email, firstname, lastname, code string) error {
	return s.mailer.SendVerificationCode(ctx, email, firstname, lastname, code)
}

// ==================================
//	Params
// ==================================
