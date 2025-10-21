package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	LogLevel     string
}

// Load configuration from environment variables
func Load() (Config, error) {
	// Try to load .env file if it exists (local development)
	_ = godotenv.Load()

	cfg := Config{
		Port:         getEnv("PORT", "8080"),
		ReadTimeout:  parseDuration(getEnv("READ_TIMEOUT", "10s")),
		WriteTimeout: parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
		IdleTimeout:  parseDuration(getEnv("IDLE_TIMEOUT", "60s")),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate configuration
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 10 * time.Second
	}
	return duration
}
