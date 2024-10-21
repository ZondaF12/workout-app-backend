package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/zondaf12/workout-app-backend/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, uuid.UUID) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, uuid.UUID)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
	}
}
