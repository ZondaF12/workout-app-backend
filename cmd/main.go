package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/zondaf12/planner-app-backend/cmd/api"
	"github.com/zondaf12/planner-app-backend/config"
	"github.com/zondaf12/planner-app-backend/db"
)

func main() {
	// Use a connection string builder function for better readability
	connStr := buildConnectionString()

	database, err := db.NewPGStorage(connStr)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}
	defer database.Close() // Ensure the database connection is closed when the program exits

	if err := initStorage(database); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Create a new Config struct for the API server
	serverConfig := api.Config{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := api.NewAPIServer(serverConfig, database)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func buildConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Envs.DBHost, config.Envs.DBPort, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName)
}

func initStorage(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database")
	return nil
}
