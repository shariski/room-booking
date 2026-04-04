package domain

import "time"

type Room struct {
	ID          string
	Name        string
	Type        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
