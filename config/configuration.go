package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	AdminUser   string
	AdminPass   string
}

func GetConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading environment variables from system")
	}

	cfg := &Config{
		Port:        mustGetEnv("PORT"),
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		AdminUser:   mustGetEnv("ADMIN_USERNAME"),
		AdminPass:   mustGetEnv("ADMIN_PASSWORD"),
	}

	return cfg
}

func mustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); !exists {
		log.Fatalf("Required environment variable %s not set", key)
	} else {
		return value
	}
	return ""
}
