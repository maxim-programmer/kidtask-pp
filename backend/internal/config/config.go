package config

import "os"

type Config struct {
	ServerAddr  string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	return Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev_secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}