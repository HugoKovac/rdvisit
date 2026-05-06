package domain

import (
	"time"

	"github.com/google/uuid"
)

type Center struct {
	ID           uuid.UUID
	Name         string
	Street       string
	StreetNumber string
	City         string
	PostalCode   string
	Region       string
	Country      string
	CreatedAt    *time.Time
	UpdatedAt    time.Time
}
