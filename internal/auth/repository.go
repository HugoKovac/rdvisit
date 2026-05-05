package auth

import (
	"context"
	"project/clean/db/sqlc"
	"project/clean/pkg/errors"
	"time"

	"github.com/google/uuid"
)

type IRepository interface {
	CreateRefreshToken(ctx context.Context, params CreateRefershToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) error
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
	// RevokeAllRefreshTokensByUser
}

// ==================================
//	Params
// ==================================

type CreateRefershToken struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}

type RevokeRefershToken struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
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

func (r *sqlcRepository) CreateRefreshToken(ctx context.Context, params CreateRefershToken) error {
	if err := r.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		ExpiresAt: params.ExpiresAt,
	}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (r *sqlcRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	refreshToken, err := r.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &RefreshToken{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: refreshToken.ExpiresAt,
		CreatedAt: refreshToken.CreatedAt,
		RevokedAt: refreshToken.RevokedAt,
	}, nil
}

func (r *sqlcRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	if err := r.q.RevokeRefreshToken(ctx, id); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (r *sqlcRepository) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteRefreshToken(ctx, id); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
