package api

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/zondaf12/workout-app-backend/service/user"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type APIServer struct {
	config Config
	db     *sql.DB
}

func NewAPIServer(config Config, db *sql.DB) *APIServer {
	return &APIServer{
		config: config,
		db:     db,
	}
}

func (s *APIServer) Start() error {
	router := fiber.New(fiber.Config{
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	})
	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	subrouter := router.Group("/api/v1")

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.config.Addr)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = router.Shutdown()
	}()

	return router.Listen(s.config.Addr)
}
