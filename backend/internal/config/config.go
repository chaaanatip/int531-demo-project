package config

import "os"

type Config struct {
	DatabaseURL    string
	ListenAddr     string
	MigrationsPath string
}

func LoadFromEnv() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// sensible default for local compose; prefer explicit env in production
		dbURL = "postgres://sreuser:srepass@database:5432/appdb?sslmode=disable"
	}
	listen := os.Getenv("LISTEN_ADDR")
	if listen == "" {
		listen = ":8080"
	}
	mig := os.Getenv("MIGRATIONS_PATH")
	if mig == "" {
		mig = "migrations"
	}
	return Config{
		DatabaseURL:    dbURL,
		ListenAddr:     listen,
		MigrationsPath: mig,
	}
}
