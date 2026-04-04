package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID
	Name        string
	Type        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
