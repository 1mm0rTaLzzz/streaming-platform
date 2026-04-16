package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	AdminUser        string
	AdminPass        string
	Env              string
	TelegramToken    string
	TelegramChatID   string
	FootballAPIKey   string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://stream:stream@localhost:5432/streamdb?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		AdminUser:      getEnv("ADMIN_USER", "admin"),
		AdminPass:      getEnv("ADMIN_PASS", "admin123"),
		Env:            getEnv("ENV", "development"),
		TelegramToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID: getEnv("TELEGRAM_CHANNEL_ID", ""),
		FootballAPIKey: getEnv("FOOTBALL_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
