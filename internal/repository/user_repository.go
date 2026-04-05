package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := "insert into users (name, email, password_hash) values ($1, $2, $3) returning id, created_at"

	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.PasswordHash).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, "Failed to insert user", "error", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := "select id, name, email, created_at from users where id = $1"

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		slog.WarnContext(ctx, "Failed to get user", "error", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := "select id, name, email, password_hash, created_at from users where email = $1"

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get user by email", "error", err)
		return nil, err
	}

	return &user, nil
}
