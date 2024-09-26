package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		GetByID(context.Context, uuid.UUID) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
	}
	Foods interface {
		Create(context.Context, *Food) error
		GetByID(context.Context, uuid.UUID) (*Food, error)
	}
	Meals interface {
		CreateMeal(context.Context, *Meal) error
		GetMeal(context.Context, Meal) (*Meal, error)
		GetMealByID(context.Context, uuid.UUID) (*Meal, error)
		CreateMealEntry(context.Context, *MealEntry) error
		GetMealEntryByID(context.Context, uuid.UUID) (*MealEntry, error)
		UpdateMealEntry(context.Context, *MealEntry) error
		DeleteMealEntry(context.Context, uuid.UUID) error
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID uuid.UUID) error
		Unfollow(ctx context.Context, followerID, userID uuid.UUID) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:     &UserStore{db},
		Foods:     &FoodStore{db},
		Meals:     &MealStore{db},
		Followers: &FollowerStore{db},
	}
}
