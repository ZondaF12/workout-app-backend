package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Meal struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Date      string    `json:"date"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type MealEntry struct {
	ID          uuid.UUID `json:"id"`
	MealID      uuid.UUID `json:"meal_id"`
	FoodID      uuid.UUID `json:"food_id"`
	ServingUnit string    `json:"serving_unit"`
	Amount      float64   `json:"amount"`
	ConsumedAt  string    `json:"consumed_at"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type MealStore struct {
	db *sql.DB
}

func (s *MealStore) CreateMeal(ctx context.Context, meal *Meal) error {
	query := `
		INSERT INTO meals (user_id, name, date)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		meal.UserID,
		meal.Name,
		meal.Date,
	).Scan(
		&meal.ID,
		&meal.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *MealStore) GetMeal(ctx context.Context, meal Meal) (*Meal, error) {
	query := `
		SELECT id, user_id, name, date, created_at, updated_at
		FROM meals
		WHERE user_id = $1 AND name = $2 AND date = $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		meal.UserID,
		meal.Name,
		meal.Date,
	).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.Name,
		&meal.Date,
		&meal.CreatedAt,
		&meal.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &meal, nil
}

func (s *MealStore) GetMealByID(ctx context.Context, id uuid.UUID) (*Meal, error) {
	query := `
		SELECT id, user_id, name, date, created_at, updated_at
		FROM meals
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var meal Meal
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.Name,
		&meal.Date,
		&meal.CreatedAt,
		&meal.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &meal, nil
}

func (s *MealStore) CreateMealEntry(ctx context.Context, entry *MealEntry) error {
	query := `
		INSERT INTO meal_entries (meal_id, food_id, serving_unit, amount, consumed_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		entry.MealID,
		entry.FoodID,
		entry.ServingUnit,
		entry.Amount,
		entry.ConsumedAt,
	).Scan(
		&entry.ID,
		&entry.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *MealStore) GetMealEntryByID(ctx context.Context, id uuid.UUID) (*MealEntry, error) {
	query := `
		SELECT id, meal_id, food_id, serving_unit, amount, consumed_at, created_at, updated_at
		FROM meal_entries
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var entry MealEntry
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.MealID,
		&entry.FoodID,
		&entry.ServingUnit,
		&entry.Amount,
		&entry.ConsumedAt,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &entry, nil
}

func (s *MealStore) UpdateMealEntry(ctx context.Context, entry *MealEntry) error {
	query := `
		UPDATE meal_entries
		SET meal_id = $1, serving_unit = $2, amount = $3, consumed_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		entry.MealID,
		entry.ServingUnit,
		entry.Amount,
		entry.ConsumedAt,
		entry.ID,
	).Scan(
		&entry.ID,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *MealStore) DeleteMealEntry(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM meal_entries
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
