package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/exec"
)

type env struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func main() {
	env, err := loadEnv()
	if err != nil {
		slog.Error("failed to load migration environment", "error", err)
		os.Exit(1)
	}

	databaseURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(env.DBUser, env.DBPassword),
		Host:     net.JoinHostPort(env.DBHost, env.DBPort),
		Path:     env.DBName,
		RawQuery: "sslmode=disable",
	}

	cmd := exec.Command("migrate", "-path", "/migrations", "-database", databaseURL.String(), "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("database migrations failed", "error", err)
		os.Exit(1)
	}
}

func loadEnv() (env, error) {
	result := env{
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	missing := make([]string, 0, 3)
	if result.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if result.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if result.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		return env{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return result, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
