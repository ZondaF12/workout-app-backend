package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", ""),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBName:                 getEnv("DB_NAME", ""),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "temporary_secret_key?"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}
	return fallback
}
