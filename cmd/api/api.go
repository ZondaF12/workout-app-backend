package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/zondaf12/workout-app-backend/docs" // This is required to load the swagger docs
	"github.com/zondaf12/workout-app-backend/internal/mailer"
	"github.com/zondaf12/workout-app-backend/internal/store"
	"go.uber.org/zap"
)

type Application struct {
	config Config
	store  store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
}

type Config struct {
	addr   string
	db     dbConfig
	env    string
	apiUrl string
	mail   mailConfig
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *Application) mount() *fiber.App {
	/* Docs */
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiUrl
	docs.SwaggerInfo.BasePath = "/v1"

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

	v1 := router.Group("/v1")

	v1.Get("/swagger/*", swagger.HandlerDefault)
	v1.Get("/health", app.healthCheckHandler)

	auth := v1.Group("/authentication")
	auth.Post("/user", app.registerUserHandler)

	users := v1.Group("/users")
	users.Put("/activate/:token", app.activateUserHandler)
	users.Get("/feed", app.getUserFeedHandler)

	user := users.Group("/:id", app.userContextMiddleware())
	user.Get("/", app.getUserHandler)
	user.Put("/follow", app.followUserHandler)
	user.Put("/unfollow", app.unfollowUserHandler)

	v1.Post("/food", app.createFoodHandler)
	v1.Get("/food/:id", app.getFoodHandler)

	v1.Post("/meal", app.createMealEntryHandler)
	v1.Patch("/meal/:id", app.updateMealEntryHandler)
	v1.Delete("/meal/:id", app.deleteMealEntryHandler)

	return router
}

func (app *Application) run(router *fiber.App) error {
	app.logger.Infow("Starting server on", "addr", app.config.addr, "env", app.config.env)

	return router.Listen(app.config.addr)
}
