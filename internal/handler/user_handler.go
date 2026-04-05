package handler

import (
	"encoding/json"
	"log/slog"
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

// @Summary Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "Register request"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Failed to decode request body", "error", err)
		writeError(w, model.NewErrBadRequest("Invalid request body"))
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		slog.WarnContext(r.Context(), "Failed to validate request body", "error", err)
		writeError(w, model.NewErrBadRequest(err.Error()))
		return
	}

	user, err := h.Usecase.Create(r.Context(), &req)
	if err != nil {
		slog.WarnContext(r.Context(), "Failed to create user", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// @Summary Login
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.LoginUserRequest true "Login request"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Failed to decode request body", "error", err)
		writeError(w, model.NewErrBadRequest("Invalid request body"))
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		slog.WarnContext(r.Context(), "Failed to validate request body", "error", err)
		writeError(w, model.NewErrBadRequest(err.Error()))
		return
	}

	auth, err := h.Usecase.Login(r.Context(), &req)
	if err != nil {
		slog.WarnContext(r.Context(), "User login failed", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, auth)
}
