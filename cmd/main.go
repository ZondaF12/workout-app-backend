package main

import (
	"log"

	"github.com/zondaf12/planner-app-backend/cmd/api"
	"github.com/zondaf12/planner-app-backend/db"
)

func main() {
	db, err := db.NewPGStorage()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
