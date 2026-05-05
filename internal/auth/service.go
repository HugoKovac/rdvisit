package auth

import (
	"context"
	"time"

	"project/clean/internal/primitive/roleprimitive"
	"project/clean/internal/user"
	"project/clean/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo   IUserRepository
	ttlAccess  time.Duration
	ttlRefresh time.Duration
	jwtSecret  string
}

type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

func NewService(userRepo IUserRepository, ttlAccess, ttlRefresh time.Duration, jwtSecret string) *Service {
	return &Service{
		userRepo:   userRepo,
		ttlAccess:  ttlAccess,
		ttlRefresh: ttlRefresh,
		jwtSecret:  jwtSecret,
	}
}

type TokenCustomClaims struct {
	jwt.RegisteredClaims
	Role roleprimitive.Role `json:"role"`
	ID   uuid.UUID          `json:"id"`
}

func (s *Service) generateToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(ttl)
	claims := &TokenCustomClaims{
		ID:   userID,
		Role: roleprimitive.Common,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.Wrap(err)
	}

	return signedToken, nil
}

func (s *Service) GenerateTokenPair(ctx context.Context, userID uuid.UUID) (string, string, error) {
	accessToken, err := s.generateToken(userID, s.ttlAccess)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateToken(userID, s.ttlRefresh)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.Wrap(errors.Unauthorized)
	}

	return s.GenerateTokenPair(ctx, user.ID)
}

func (s *Service) ValidateRefreshToken(tokenString string) (*TokenCustomClaims, error) {
	claims := &TokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, errors.Wrap(errors.Unauthorized)
	}

	if !token.Valid {
		return nil, errors.Wrap(errors.Unauthorized)
	}

	claims, ok := token.Claims.(*TokenCustomClaims)
	if !ok {
		return nil, errors.Wrap(errors.Unauthorized)
	}

	return claims, nil
}

// ==================================
//	Params
// ==================================

