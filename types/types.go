package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	CreateUser(User) error
}

type FoodStore interface {
	GetAllFoods() (*[]Food, error)
	CreateFood(Food) error
}

type User struct {
	ID        uuid.UUID      `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Bio       sql.NullString `json:"bio"`
	CreatedAt time.Time      `json:"created_at"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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
