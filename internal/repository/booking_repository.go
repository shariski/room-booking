package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/domain"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking domain.Booking) (*domain.Booking, error) {
	query := "insert into bookings (room_id, user_id, start_date, end_date) values ($1, $2, $3, $4) returning id, created_at, updated_at"

	err := r.db.QueryRowContext(ctx, query, booking.RoomID, booking.UserID, booking.StartDate, booking.EndDate).Scan(
		&booking.ID,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, "Failed to insert booking", "error", err)
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Booking, error) {
	query := "update bookings SET deleted_at = NOW() where id = $1 and user_id = $2 and deleted_at is null returning id, room_id, user_id, start_date, end_date, created_at, updated_at, deleted_at"

	var booking domain.Booking
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&booking.ID,
		&booking.RoomID,
		&booking.UserID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.DeletedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, "Failed to delete booking", "error", err)
		return nil, err
	}

	return &booking, nil
}
