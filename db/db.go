package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zondaf12/planner-app-backend/config"
)

func NewPGStorage() (*pgx.Conn, error) {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBAddress,
		config.Envs.DBName,
	)

	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	log.Println("Connected to database")

	return conn, nil
}
