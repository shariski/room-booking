package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/model"
	"github.com/shariski/room-booking/internal/usecase"
)

type BookingHandler struct {
	Usecase  *usecase.BookingUsecase
	Validate *validator.Validate
}

func NewBookingHandler(u *usecase.BookingUsecase, v *validator.Validate) *BookingHandler {
	return &BookingHandler{
		Usecase:  u,
		Validate: v,
	}
}

// @Summary Create a booking
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateBookingRequest true "Booking request"
// @Success 201 {object} model.BookingResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /bookings [post]
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Failed to decode request body", "error", err)
		writeError(w, model.NewErrBadRequest("Invalid request body"))
		return
	}

	claims := r.Context().Value("user").(*model.Auth)
	req.UserID = claims.ID

	if err := h.Validate.Struct(req); err != nil {
		slog.WarnContext(r.Context(), "Failed to validate request body", "error", err)
		writeError(w, model.NewErrBadRequest(err.Error()))
		return
	}

	booking, err := h.Usecase.Create(r.Context(), &req)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed to create booking", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, booking)
}

// @Summary Cancel a booking
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID"
// @Success 200 {object} model.BookingResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /bookings/{id} [delete]
func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	paramID := r.PathValue("id")
	bookingID, err := uuid.Parse(paramID)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed to parse UUID", "error", err)
		writeError(w, model.NewErrBadRequest("Invalid UUID"))
		return
	}

	claims := r.Context().Value("user").(*model.Auth)
	req := &model.DeleteBookingRequest{
		ID:     bookingID,
		UserID: claims.ID,
	}

	if err := h.Validate.Struct(req); err != nil {
		slog.WarnContext(r.Context(), "Failed to validate request", "error", err)
		writeError(w, model.NewErrBadRequest(err.Error()))
		return
	}

	booking, err := h.Usecase.Delete(r.Context(), req)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed to delete booking", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, booking)
}
