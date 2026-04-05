package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shariski/room-booking/internal/config"
	"github.com/shariski/room-booking/internal/domain"
	"github.com/shariski/room-booking/internal/model"
	"github.com/shariski/room-booking/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Config *config.Config
	Repo   *repository.UserRepository
}

func NewUserUsecase(config *config.Config, repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{Config: config, Repo: repo}
}

func (u *UserUsecase) Create(ctx context.Context, request *model.CreateUserRequest) (*model.UserResponse, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.WarnContext(ctx, "Failed to generate bcrypt hash", "error", err)
		return nil, err
	}

	userData := &domain.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(password),
	}

	user, err := u.Repo.Create(ctx, userData)
	if err != nil {
		var dbErr *pgconn.PgError
		if errors.As(err, &dbErr) && dbErr.Code == "23505" {
			return nil, model.NewConflictError("Email already registered")
		}
		slog.WarnContext(ctx, "Failed to create user", "error", err)
		return nil, err
	}

	return model.UserToResponse(user), nil
}

func (u *UserUsecase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.AuthResponse, error) {
	user, err := u.Repo.GetByEmail(ctx, request.Email)
	if errors.Is(err, sql.ErrNoRows) {
		slog.WarnContext(ctx, "User not found", "error", err)
		return nil, model.NewErrBadRequest("Incorrect credential")
	}
	if err != nil {
		slog.WarnContext(ctx, "Failed to find user by email", "error", err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		slog.WarnContext(ctx, "Failed to compare password", "error", err)
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})

	secret := u.Config.JWTSecret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.WarnContext(ctx, "Failed to sign jwt token", "error", err)
		return nil, err
	}

	return &model.AuthResponse{ID: user.ID, AccessToken: tokenString}, nil
}
