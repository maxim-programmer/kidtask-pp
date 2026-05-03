package config

import "os"

type Config struct {
	ServerAddr  string
	DatabaseURL string
	JWTSecret   string
	AdminSecret string
}

func Load() Config {
	return Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev_secret"),
		AdminSecret: getEnv("ADMIN_SECRET", "admin-secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}