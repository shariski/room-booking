package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/domain"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) List(ctx context.Context, roomType string) ([]domain.Room, error) {
	query := "select id, name, type, description, created_at, updated_at from rooms"
	args := []any{}

	if roomType != "" {
		query += " where type = $1"
		args = append(args, roomType)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.WarnContext(ctx, "Failed to do query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Type,
			&room.Description,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			slog.WarnContext(ctx, "Failed to scan data", "error", err)
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, rows.Err()
}

func (r *RoomRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	query := "select id, name, type, description, created_at, updated_at from rooms where id = $1"

	var room domain.Room
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
		&room.Type,
		&room.Description,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, "Failed to do query", "error", err)
		return nil, err
	}

	return &room, nil
}
