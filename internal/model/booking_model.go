package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/domain"
)

type BookingResponse struct {
	ID        uuid.UUID  `json:"id"`
	RoomID    uuid.UUID  `json:"room_id"`
	UserID    uuid.UUID  `json:"user_id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateBookingRequest struct {
	RoomID    uuid.UUID `json:"room_id" validate:"required"`
	UserID    uuid.UUID `json:"-"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

type DeleteBookingRequest struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	UserID uuid.UUID `json:"-"`
}

func BookingToResponse(booking *domain.Booking) *BookingResponse {
	return &BookingResponse{
		ID:        booking.ID,
		RoomID:    booking.RoomID,
		UserID:    booking.UserID,
		StartDate: booking.StartDate,
		EndDate:   booking.EndDate,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
		DeletedAt: booking.DeletedAt,
	}
}
