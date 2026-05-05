package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"project/clean/internal/domain"
	"project/clean/internal/primitive/roleprimitive"
	"project/clean/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo                IRepository
	userRepo            IUserRepository
	mailer              IMailer
	ttlAccess           time.Duration
	ttlRefresh          time.Duration
	jwtSecret           string
	verificationCodeTTL time.Duration
}

type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, email, firstname, lastname, passwordHash string) (*domain.User, error)
	VerifyUserByID(ctx context.Context, id uuid.UUID) error
}

func NewService(repo IRepository, userRepo IUserRepository, mailer IMailer, ttlAccess, ttlRefresh, verificationCodeTTL time.Duration, jwtSecret string) *Service {
	return &Service{
		repo:                repo,
		userRepo:            userRepo,
		mailer:              mailer,
		ttlAccess:           ttlAccess,
		ttlRefresh:          ttlRefresh,
		jwtSecret:           jwtSecret,
		verificationCodeTTL: verificationCodeTTL,
	}
}

type TokenCustomClaims struct {
	jwt.RegisteredClaims
	Role roleprimitive.Role `json:"role"`
	ID   uuid.UUID          `json:"id"`
}

func generateCode() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *Service) Verify(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	vc, err := s.repo.GetVerificationCodeByUser(ctx, userID)
	if err != nil {
		return false, err
	}
	if strings.Compare(vc.Code, code) != 0 && vc.ExpiresAt.After(time.Now()) {
		return false, errors.Wrap(errors.Unauthorized)
	}
	if err := s.userRepo.VerifyUserByID(ctx, userID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Register(ctx context.Context, email, firstname, lastname, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	//todo: params validation and normalisation

	user, err := s.userRepo.Create(ctx, email, firstname, lastname, string(hash))
	if err != nil {
		return nil, err
	}

	code, err := generateCode()
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateVerificationCode(ctx, user.ID, code, time.Now().Add(s.verificationCodeTTL)); err != nil {
		return nil, err
	}

	if err := s.mailer.SendVerificationCode(ctx, email, firstname, lastname, code); err != nil {
		return nil, err
	}

	return user, nil
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

func (s *Service) DeleteRefreshToken(ctx context.Context, tokenString string) error {
	hash := sha256.Sum256([]byte(tokenString))
	refreshToken, err := s.repo.GetRefreshTokenByHash(ctx, hex.EncodeToString(hash[:]))
	if err != nil {
		return errors.Wrap(errors.Unauthorized)
	}
	if refreshToken.RevokedAt != nil && refreshToken.RevokedAt.Before(time.Now()) {
		return errors.Wrap(errors.Unauthorized)
	}
	if err := s.repo.DeleteRefreshToken(ctx, refreshToken.ID); err != nil {
		return errors.Wrap(errors.Unauthorized)
	}
	return nil
}

func (s *Service) ValidateToken(tokenString string) (*TokenCustomClaims, error) {
	claims := &TokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err)
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
