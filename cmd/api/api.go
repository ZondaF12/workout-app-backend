package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/zondaf12/planner-app-backend/service/user"
)

type APIServer struct {
	addr string
	db   *pgx.Conn
}

func NewAPIServer(addr string, db *pgx.Conn) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Start() error {
	router := fiber.New()
	subrouter := router.Group("/api/v1")

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)

	return router.Listen(s.addr)
}
