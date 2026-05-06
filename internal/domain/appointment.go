package domain

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID             uuid.UUID
	Date           time.Time
	PractitionerID uuid.UUID
	PatientEmail   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
