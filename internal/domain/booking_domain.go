package domain

import "time"

type Booking struct {
	ID        string
	RoomID    string
	UserID    string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
