package domain

import (
	"time"

	"github.com/google/uuid"
)

type VerificationCode struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Code       string
	ExpiresAt  time.Time
	Created_at time.Time
}
