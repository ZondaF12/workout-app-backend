package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/zondaf12/workout-app-backend/internal/db"
	"github.com/zondaf12/workout-app-backend/internal/env"
	"github.com/zondaf12/workout-app-backend/internal/store"
	"go.uber.org/zap"

	_ "github.com/zondaf12/workout-app-backend/docs"
)

const version = "0.0.1"

// @title						Workout App API
// @description				This is the API documentation for the Workout App API.
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.email				fiber@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	godotenv.Load()

	cfg := Config{
		addr:   env.GetString("ADDR", ":8080"),
		apiUrl: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost:5432/workoutapp?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Main Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)

	app := &Application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	router := app.mount()
	logger.Fatal(app.run(router))
}
