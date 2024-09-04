package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/zondaf12/planner-app-backend/cmd/api"
	"github.com/zondaf12/planner-app-backend/config"
	"github.com/zondaf12/planner-app-backend/db"
)

func main() {
	db, err := db.NewPGStorage(fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		config.Envs.DBHost, config.Envs.DBPort, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName))
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
}
