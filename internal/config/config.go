package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")

	jwtSecret := os.Getenv("JWT_SECRET")

	return Config{Port: port, DatabaseURL: databaseURL, JWTSecret: jwtSecret}
}
