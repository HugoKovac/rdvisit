package user

import (
	"context"
	"project/clean/db/sqlc"
	"project/clean/pkg/errors"
)

type IRepository interface {
	Create(ctx context.Context, params CreateParams) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
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

func (r *sqlcRepository) Create(ctx context.Context, params CreateParams) (*User, error) {
	user, err := r.q.Create(ctx, sqlc.CreateParams{
		Email:        params.Email,
		Firstname:    params.FirstName,
		Lastname:     params.LastName,
		PasswordHash: params.PasswordHash,
	})
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &User{
		ID:           user.ID,
		FirstName:    user.Firstname,
		LastName:     user.Lastname,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (r *sqlcRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	user, err := r.q.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &User{
		ID:           user.ID,
		FirstName:    user.Firstname,
		LastName:     user.Lastname,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
