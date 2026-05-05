package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"project/clean/internal/primitive/roleprimitive"
	"project/clean/internal/user"
	"project/clean/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo       IRepository
	userRepo   IUserRepository
	ttlAccess  time.Duration
	ttlRefresh time.Duration
	jwtSecret  string
}

type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

func NewService(repo IRepository, userRepo IUserRepository, ttlAccess, ttlRefresh time.Duration, jwtSecret string) *Service {
	return &Service{
		repo:       repo,
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

func (s *Service) generateToken(userID uuid.UUID) (string, error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(s.ttlAccess)
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

func (s *Service) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(s.ttlRefresh)
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

	hash := sha256.Sum256([]byte(signedToken))
	if err := s.repo.CreateRefreshToken(ctx, CreateRefershToken{
		UserID:    userID,
		TokenHash: hex.EncodeToString(hash[:]),
		ExpiresAt: expirationTime,
	}); err != nil {
		return "", errors.Wrap(err)
	}

	return signedToken, nil
}

func (s *Service) GenerateTokenPair(ctx context.Context, userID uuid.UUID) (string, string, error) {
	accessToken, err := s.generateToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateRefreshToken(ctx, userID)
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

func (s *Service) RevokeRefreshToken(ctx context.Context, tokenString string) error {
	hash := sha256.Sum256([]byte(tokenString))
	refreshToken, err := s.repo.GetRefreshTokenByHash(ctx, hex.EncodeToString(hash[:]))
	if err != nil {
		return errors.Wrap(err)
	}
	if refreshToken.RevokedAt != nil && refreshToken.RevokedAt.Before(time.Now()) {
		return errors.Wrap(errors.Unauthorized)
	}
	if err := s.repo.RevokeRefreshToken(ctx, refreshToken.ID); err != nil {
		return errors.Wrap(err)
	}
	return nil
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
