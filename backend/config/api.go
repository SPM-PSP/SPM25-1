package config

import (
	"os"
	"time"
)

type AppConfig struct {
	DeepSeekAPIKey string
	Port           string
	APITimeout     time.Duration
}

func Load() *AppConfig {
	return &AppConfig{
		DeepSeekAPIKey: os.Getenv("DEEPSEEK_API_KEY"),
		Port:           getEnv("PORT", "8080"),
		APITimeout:     30 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
