package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shariski/room-booking/internal/model"
	"github.com/shariski/room-booking/internal/usecase"
)

type UserHandler struct {
	Usecase  *usecase.UserUsecase
	Validate *validator.Validate
}

func NewUserHandler(u *usecase.UserUsecase, v *validator.Validate) *UserHandler {
	return &UserHandler{
		Usecase:  u,
		Validate: v,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.Create(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	auth, err := h.Usecase.Login(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(auth)
}
