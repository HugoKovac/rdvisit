package user

import (
	"context"
	"project/clean/db/sqlc"
	"project/clean/internal/domain"
	"project/clean/pkg/errors"
	"time"

	"github.com/google/uuid"
)

type IRepository interface {
	Create(ctx context.Context, email, firstname, lastname, passwordHash string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	VerifyUserByID(ctx context.Context, id uuid.UUID) error

	CreateVerificationCode(ctx context.Context, userID uuid.UUID, code string, expiresAt time.Time) error
	GetVerificationCodeByUser(ctx context.Context, userID uuid.UUID) (*domain.VerificationCode, error)
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

func (r *sqlcRepository) Create(ctx context.Context, email, firstname, lastname, passwordHash string) (*domain.User, error) {
	user, err := r.q.Create(ctx, sqlc.CreateParams{
		Email:        email,
		Firstname:    firstname,
		Lastname:     lastname,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &domain.User{
		ID:           user.ID,
		FirstName:    user.Firstname,
		LastName:     user.Lastname,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (r *sqlcRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.q.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &domain.User{
		ID:           user.ID,
		FirstName:    user.Firstname,
		LastName:     user.Lastname,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (r *sqlcRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := r.q.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &domain.User{
		ID:           user.ID,
		FirstName:    user.Firstname,
		LastName:     user.Lastname,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (r *sqlcRepository) VerifyUserByID(ctx context.Context, id uuid.UUID) error {
	if err := r.q.VerifyUserByID(ctx, id); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (r *sqlcRepository) CreateVerificationCode(ctx context.Context, userID uuid.UUID, code string, expiresAt time.Time) error {
	err := r.q.CreateVerificationCode(ctx, sqlc.CreateVerificationCodeParams{
		UserID:    userID,
		Code:      code,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (r *sqlcRepository) GetVerificationCodeByUser(ctx context.Context, userID uuid.UUID) (*domain.VerificationCode, error) {
	vc, err := r.q.GetVerificationCode(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &domain.VerificationCode{
		ID:         vc.ID,
		UserID:     vc.UserID,
		Code:       vc.Code,
		ExpiresAt:  vc.ExpiresAt,
		Created_at: vc.CreatedAt,
	}, nil
}
