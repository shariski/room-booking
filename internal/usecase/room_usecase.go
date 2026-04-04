package usecase

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/shariski/room-booking/internal/model"
	"github.com/shariski/room-booking/internal/repository"
)

type RoomUsecase struct {
	Repo  *repository.RoomRepository
	Redis *redis.Client
}

func NewRoomUsecase(repo *repository.RoomRepository, redis *redis.Client) *RoomUsecase {
	return &RoomUsecase{Repo: repo, Redis: redis}
}

func (u *RoomUsecase) List(ctx context.Context, request *model.GetRoomsRequest) ([]model.RoomResponse, error) {
	// TODO: get cache from redis
	rooms, err := u.Repo.List(ctx, request.Type)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get room list", "error", err)
		return nil, err
	}

	return model.RoomsToResponse(rooms), nil
}

func (u *RoomUsecase) Get(ctx context.Context, request *model.GetRoomRequest) (*model.RoomResponse, error) {
	room, err := u.Repo.Get(ctx, request.ID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get room", "error", err)
		return nil, err
	}

	return model.RoomToResponse(room), nil
}
