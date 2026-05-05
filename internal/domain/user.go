package domain

import (
	"project/clean/internal/primitive/roleprimitive"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	Role         roleprimitive.Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
