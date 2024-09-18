package types

import (
	"time"

	"github.com/google/uuid"
)

type FoodStore interface {
	GetAllFoods() (*[]Food, error)
	CreateFood(Food) error
}

type Food struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	Brand              *string    `json:"brand"`
	DefaultServingSize float64    `json:"default_serving_size"`
	DefaultServingUnit string     `json:"default_serving_unit"`
	Calories           int        `json:"calories"`
	Protein            *float64   `json:"protein"`
	Carbs              *float64   `json:"carbs"`
	Fat                *float64   `json:"fat"`
	IsUserCreated      bool       `json:"is_user_created"`
	CreatedBy          *uuid.UUID `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
}
