package handler

import (
	"log/slog"
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

// @Summary List rooms
// @Tags rooms
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter by room type (e.g. single, double, suite)"
// @Success 200 {array} model.RoomResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /rooms [get]
func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	qType := query.Get("type")

	req := &model.GetRoomsRequest{Type: qType}
	if err := h.Validate.Struct(req); err != nil {
		slog.WarnContext(r.Context(), "Failed to validate request", "error", err)
		writeError(w, model.NewErrBadRequest(err.Error()))
		return
	}

	rooms, err := h.Usecase.List(r.Context(), req)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed to get list room", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, rooms)
}

// @Summary Get room detail
// @Tags rooms
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 200 {object} model.RoomResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /rooms/{id} [get]
func (h *RoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	paramID := r.PathValue("id")
	roomID, err := uuid.Parse(paramID)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed to parse UUID", "error", err)
		writeError(w, model.NewErrBadRequest("Invalid UUID"))
		return
	}

	req := &model.GetRoomRequest{ID: roomID}

	if err := h.Validate.Struct(req); err != nil {
		slog.WarnContext(r.Context(), "Failed to validate request", "error", err)
		writeError(w, model.NewErrBadRequest(err.Error()))
		return
	}

	room, err := h.Usecase.Get(r.Context(), req)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed to get room", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, room)
}
