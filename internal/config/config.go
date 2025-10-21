package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DefaultPackSizes []int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	LogLevel         string
}

// Load configuration from environment variables
func Load() (Config, error) {
	// Try to load .env file if it exists (local development)
	_ = godotenv.Load()

	cfg := Config{
		Port:             getEnv("PORT", "8080"),
		DefaultPackSizes: parsePackSizes(getEnv("DEFAULT_PACK_SIZES", "250,500,1000,2000,5000")),
		ReadTimeout:      parseDuration(getEnv("READ_TIMEOUT", "10s")),
		WriteTimeout:     parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
		IdleTimeout:      parseDuration(getEnv("IDLE_TIMEOUT", "60s")),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
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

	if len(c.DefaultPackSizes) == 0 {
		return fmt.Errorf("DEFAULT_PACK_SIZES cannot be empty")
	}

	for _, size := range c.DefaultPackSizes {
		if size <= 0 {
			return fmt.Errorf("pack sizes must be positive, got: %d", size)
		}
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parsePackSizes(value string) []int {
	parts := strings.Split(value, ",")
	sizes := make([]int, 0, len(parts))

	for _, part := range parts {
		if size, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			sizes = append(sizes, size)
		}
	}

	return sizes
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 10 * time.Second
	}
	return duration
}
