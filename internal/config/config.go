package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 5*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
