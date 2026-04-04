package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/model"
	"github.com/shariski/room-booking/internal/usecase"
)

type RoomHandler struct {
	Usecase  *usecase.RoomUsecase
	Validate *validator.Validate
}

func NewRoomHandler(u *usecase.RoomUsecase, v *validator.Validate) *RoomHandler {
	return &RoomHandler{
		Usecase:  u,
		Validate: v,
	}
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	qType := query.Get("type")

	req := &model.GetRoomsRequest{Type: qType}
	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	rooms, err := h.Usecase.List(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	paramID := r.PathValue("id")
	roomID, err := uuid.Parse(paramID)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	req := &model.GetRoomRequest{ID: roomID}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	room, err := h.Usecase.Get(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}
