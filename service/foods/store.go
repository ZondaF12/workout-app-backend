package foods

import (
	"database/sql"
	"fmt"

	"github.com/zondaf12/workout-app-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetAllFoods() (*[]types.Food, error) {
	query := "SELECT * FROM foods"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching foods: %w", err)
	}
	defer rows.Close()

	foods := []types.Food{}

	for rows.Next() {
		var food types.Food
		err := rows.Scan(
			&food.ID, &food.Name, &food.Brand, &food.DefaultServingSize,
			&food.DefaultServingUnit, &food.Calories, &food.Protein,
			&food.Carbs, &food.Fat, &food.IsUserCreated, &food.CreatedBy, &food.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning food: %w", err)
		}
		foods = append(foods, food)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching foods: %w", err)
	}

	return &foods, nil
}

func (s *Store) CreateFood(food types.Food) error {
	return nil
}
