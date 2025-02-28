package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		DSN string
	}
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Server configuration
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// Database configuration
	cfg.Database.DSN = getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/scheduler?sslmode=disable")

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt gets an environment variable as an integer or returns a default value --- Might not need it
func getEnvInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
