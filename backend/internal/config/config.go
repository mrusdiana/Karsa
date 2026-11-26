package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBDSN       string
	JWTSecret   string
	FrontendURL string
	APIURL      string
	UploadDir   string
	SMTPHost    string
	SMTPPort    string
	MailFrom    string
	RunSeed     bool
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:        get("PORT", "8080"),
		DBDSN:       get("DB_DSN", ""),
		JWTSecret:   get("JWT_SECRET", ""),
		FrontendURL: get("FRONTEND_URL", "http://localhost:5173"),
		APIURL:      get("API_URL", "http://localhost:8080"),
		UploadDir:   get("UPLOAD_DIR", "uploads"),
		SMTPHost:    get("SMTP_HOST", "localhost"),
		SMTPPort:    get("SMTP_PORT", "1025"),
		MailFrom:    get("MAIL_FROM", "no-reply@karsa.test"),
		RunSeed:     get("RUN_SEED", "false") == "true",
	}
}

func get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
