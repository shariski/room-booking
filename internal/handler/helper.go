package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shariski/room-booking/internal/model"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		writeJSON(w, appErr.Status, model.ErrorResponse{
			Message: appErr.Message,
		})
		return
	}

	writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
		Message: "internal server error",
	})
}
