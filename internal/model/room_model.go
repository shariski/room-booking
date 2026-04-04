package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shariski/room-booking/internal/domain"
)

type RoomResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetRoomsRequest struct {
	Type string `json:"type"`
}

type GetRoomRequest struct {
	ID uuid.UUID `json:"id"`
}

func RoomToResponse(room *domain.Room) *RoomResponse {
	return &RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Type:        room.Type,
		Description: room.Description,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}
}

func RoomsToResponse(rooms []domain.Room) []RoomResponse {
	res := make([]RoomResponse, len(rooms))

	for i, r := range rooms {
		res[i] = RoomResponse{
			ID:          r.ID,
			Name:        r.Name,
			Type:        r.Type,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}

	return res
}
