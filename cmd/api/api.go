package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/zondaf12/workout-app-backend/docs" // This is required to load the swagger docs
	"github.com/zondaf12/workout-app-backend/internal/auth"
	"github.com/zondaf12/workout-app-backend/internal/mailer"
	"github.com/zondaf12/workout-app-backend/internal/ratelimiter"
	"github.com/zondaf12/workout-app-backend/internal/store"
	"github.com/zondaf12/workout-app-backend/internal/store/cache"
	"go.uber.org/zap"
)

type Application struct {
	config        Config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type Config struct {
	addr        string
	db          dbConfig
	env         string
	apiUrl      string
	mail        mailConfig
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	basic basicAuthConfig
	token tokenConfig
}

type basicAuthConfig struct {
	username string
	password string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
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
	router.Use(app.RateLimiterMiddleware())

	v1 := router.Group("/v1")

	v1.Get("/swagger/*", swagger.HandlerDefault)
	v1.Get("/health", app.healthCheckHandler)

	auth := v1.Group("/authentication")
	auth.Post("/register", app.registerUserHandler)
	auth.Post("/login", app.BasicAuthMiddleware(), app.createTokenHandler)

	users := v1.Group("/users")
	users.Put("/activate/:token", app.activateUserHandler)
	users.Get("/self", app.AuthTokenMiddleware(), app.getSelfHandler)
	users.Get("/feed", app.AuthTokenMiddleware(), app.getUserFeedHandler)

	user := users.Group("/:id", app.AuthTokenMiddleware())
	user.Get("/", app.AuthTokenMiddleware(), app.getUserHandler)
	user.Put("/follow", app.AuthTokenMiddleware(), app.followUserHandler)
	user.Put("/unfollow", app.AuthTokenMiddleware(), app.unfollowUserHandler)

	v1.Post("/food", app.AuthTokenMiddleware(), app.createFoodHandler)
	v1.Get("/food", app.AuthTokenMiddleware(), app.getAllFoodHandler)
	v1.Get("/food/:id", app.AuthTokenMiddleware(), app.getFoodHandler)

	meals := v1.Group("/meals", app.AuthTokenMiddleware())
	meals.Post("/", app.createMealEntryHandler)
	meals.Patch("/:id", app.updateMealEntryHandler)
	meals.Delete("/:id", app.deleteMealEntryHandler)

	return router
}

func (app *Application) run(router *fiber.App) error {
	// Channel for shutdown errors
	shutdown := make(chan error, 1)

	// Graceful shutdown goroutine
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Infow("signal caught", "signal", s.String())

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown with context
		if err := router.ShutdownWithContext(ctx); err != nil {
			app.logger.Errorw("shutdown error", "error", err)
			shutdown <- err
			return
		}

		shutdown <- nil
	}()

	app.logger.Infow("Starting server on", "addr", app.config.addr, "env", app.config.env)

	err := router.Listen(app.config.addr)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
