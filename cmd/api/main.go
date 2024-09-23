package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/zondaf12/workout-app-backend/internal/db"
	"github.com/zondaf12/workout-app-backend/internal/env"
	"github.com/zondaf12/workout-app-backend/internal/store"

	_ "github.com/zondaf12/workout-app-backend/docs"
)

//	@title						Workout App API
//	@version					1.0
//	@description				This is the API documentation for the Workout App API.
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.email				fiber@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8080
//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
func main() {
	godotenv.Load()

	cfg := Config{
		Addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost:5432/workoutapp?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	// Main Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(nil)

	app := &Application{
		config: cfg,
		store:  store,
	}

	router := app.mount()
	log.Fatal(app.run(router))
}
