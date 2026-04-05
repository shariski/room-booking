package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

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
	rooms, err := u.Repo.List(ctx, request.Type)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get room list", "error", err)
		return nil, err
	}

	return model.RoomsToResponse(rooms), nil
}

func (u *RoomUsecase) Get(ctx context.Context, request *model.GetRoomRequest) (*model.RoomResponse, error) {
	cacheKey := "room:" + request.ID.String()

	cached, err := u.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var room model.RoomResponse
		json.Unmarshal([]byte(cached), &room)
		return &room, nil
	} else {
		// no return error, just fallback to DB
		slog.WarnContext(ctx, "Failed to get redis cache", "error", err)
	}

	room, err := u.Repo.Get(ctx, request.ID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get room", "error", err)
		return nil, err
	}

	response := model.RoomToResponse(room)

	data, _ := json.Marshal(response)
	u.Redis.Set(ctx, cacheKey, data, 10*time.Minute)

	return response, nil
}
