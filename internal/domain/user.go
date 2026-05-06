package domain

import (
	"time"

	"github.com/HugoKovac/rdvisit/internal/primitive/roleprimitive"

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
