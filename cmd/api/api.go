package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/zondaf12/workout-app-backend/internal/store"

	_ "github.com/zondaf12/workout-app-backend/docs"
)

type Application struct {
	config Config
	store  store.Storage
}

type Config struct {
	Addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *Application) mount() *fiber.App {
	srv := fiber.Config{
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	router := fiber.New(srv)

	router.Use(recover.New())
	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	router.Get("/swagger/*", swagger.HandlerDefault)

	version := router.Group("/v1")

	version.Get("/health", app.healthCheckHandler)

	return router
}

func (app *Application) run(router *fiber.App) error {
	log.Println("Starting server on", app.config.Addr)

	return router.Listen(app.config.Addr)
}
