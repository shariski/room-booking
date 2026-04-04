package handler

import (
	"encoding/json"
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

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	booking, err := h.Usecase.Create(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	paramID := r.PathValue("id")
	bookingID, err := uuid.Parse(paramID)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	req := &model.DeleteBookingRequest{ID: bookingID}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	booking, err := h.Usecase.Delete(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}
