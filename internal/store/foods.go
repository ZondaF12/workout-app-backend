package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Food struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Calories    int       `json:"calories"`
	Protein     float64   `json:"protein"`
	Carbs       float64   `json:"carbs"`
	Fat         float64   `json:"fat"`
	Brand       string    `json:"brand"`
	ServingSize float64   `json:"serving_size"`
	ServingUnit string    `json:"serving_unit"`
	Verified    bool      `json:"verified"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type FoodStore struct {
	db *sql.DB
}

func (s *FoodStore) Create(ctx context.Context, food *Food) error {
	query := `
		INSERT INTO foods (name, description, calories, protein, carbs, fat, brand, serving_size, serving_unit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		food.Name,
		food.Description,
		food.Calories,
		food.Protein,
		food.Carbs,
		food.Fat,
		food.Brand,
		food.ServingSize,
		food.ServingUnit,
	).Scan(
		&food.ID,
		&food.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *FoodStore) GetByID(ctx context.Context, id uuid.UUID) (*Food, error) {
	query := `
		SELECT id, name, description, calories, protein, carbs, fat, brand, serving_size, serving_unit, verified, user_id, created_at, updated_at
		FROM foods
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var food Food
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&food.ID,
		&food.Name,
		&food.Description,
		&food.Calories,
		&food.Protein,
		&food.Carbs,
		&food.Fat,
		&food.Brand,
		&food.ServingSize,
		&food.ServingUnit,
		&food.Verified,
		&food.UserID,
		&food.CreatedAt,
		&food.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &food, nil
}
