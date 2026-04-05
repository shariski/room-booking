package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shariski/room-booking/internal/domain"
	"github.com/shariski/room-booking/internal/model"
	"github.com/shariski/room-booking/internal/repository"
)

type BookingUsecase struct {
	repo *repository.BookingRepository
}

func NewBookingUsecase(repo *repository.BookingRepository) *BookingUsecase {
	return &BookingUsecase{repo: repo}
}

func (u *BookingUsecase) Create(ctx context.Context, request *model.CreateBookingRequest) (*model.BookingResponse, error) {
	bookingData := &domain.Booking{
		RoomID:    request.RoomID,
		UserID:    request.UserID,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	}
	booking, err := u.repo.Create(ctx, *bookingData)
	if err != nil {
		var dbErr *pgconn.PgError
		// 23P01: exclusion constraint violation — overlapping booking dates
		if errors.As(err, &dbErr) && dbErr.Code == "23P01" {
			slog.WarnContext(ctx, "Booking conflict with overlapping bookings", "error", err)
			return nil, model.NewConflictError("Booking conflict with overlapping dates")
			// 23514: check constraint - end date should be larger than start date
		} else if errors.As(err, &dbErr) && dbErr.Code == "23514" {
			slog.WarnContext(ctx, "Booking start date cannot larger than end date", "error", err)
			return nil, model.NewErrBadRequest("Booking start date cannot larger than end date")
		}
		slog.WarnContext(ctx, "Failed to create booking", "error", err)
		return nil, err
	}

	return model.BookingToResponse(booking), nil
}

func (u *BookingUsecase) Delete(ctx context.Context, request *model.DeleteBookingRequest) (*model.BookingResponse, error) {
	booking, err := u.repo.Delete(ctx, request.ID, request.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		slog.WarnContext(ctx, "Booking not found", "error", err)
		return nil, model.NewErrNotFound("Booking not found")
	}
	if err != nil {
		slog.WarnContext(ctx, "Failed to delete booking", "error", err)
		return nil, err
	}

	return model.BookingToResponse(booking), nil
}
