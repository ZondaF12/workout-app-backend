package user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/zondaf12/planner-app-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := s.db.QueryRow(query, email)

	u := new(types.User)
	err := scanRowIntoUser(row, u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return u, nil
}

func scanRowIntoUser(row *sql.Row, user *types.User) error {
	return row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.CreatedAt,
	)
}

func (s *Store) GetUserByID(id uuid.UUID) (*types.User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	row := s.db.QueryRow(query, id)

	u := new(types.User)
	err := scanRowIntoUser(row, u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return u, nil
}

func (s *Store) CreateUser(u types.User) error {
	query := "INSERT INTO users (id, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(query, uuid.New(), u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}
