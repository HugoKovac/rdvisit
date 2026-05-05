package user

import (
	"context"
	"project/clean/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, params RegisterParams) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	//todo: params validation and normalisation

	return s.repo.Create(ctx, CreateParams{
		params.Email,
		params.FirstName,
		params.LastName,
		string(hash),
	})
}

// ==================================
//	Params
// ==================================

type RegisterParams struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}
