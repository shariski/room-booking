package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
