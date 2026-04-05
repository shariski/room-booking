package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/domain"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResponse struct {
	ID          uuid.UUID `json:"id"`
	AccessToken string    `json:"access_token"`
}

type Auth struct {
	ID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyUserRequest struct {
	AccessToken string
}

func UserToResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
