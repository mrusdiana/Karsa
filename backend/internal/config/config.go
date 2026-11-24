package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBDSN       string
	JWTSECRET   string
	FrontendURL string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:        get("PORT", "8080"),
		DBDSN:       get("DB_DSN", ""),
		JWTSECRET:   get("JWT_SECRET", ""),
		FrontendURL: get("FRONTEND_URL", "http://localhost:5173"),
	}
}

func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
