package main

import (
	"log"

	"github.com/zondaf12/planner-app-backend/cmd/api"
)

func main() {
	server := api.NewAPIServer(":8080", nil)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
