package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	PostgresDSN string
	HTTPPort    string
}

func MustLoad() *Config {
	requiredVars := []string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"HTTP_PORT",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("❌ missing required environment variable: %s", v)
		}
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	return &Config{
		PostgresDSN: dsn,
		HTTPPort:    os.Getenv("HTTP_PORT"),
	}
}
