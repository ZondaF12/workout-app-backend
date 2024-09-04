package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432")),
		DBName:     getEnv("DB_NAME", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
