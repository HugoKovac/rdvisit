package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"project/clean/db/sqlc"
	"project/clean/internal/auth"
	"project/clean/internal/user"
	loggerMiddleware "project/clean/pkg/logger/middleware"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	APP_PORT        string        `envconfig:"APP_PORT" required:"true"`
	DB_HOST         string        `envconfig:"DB_HOST" required:"true"`
	DB_PORT         string        `envconfig:"DB_PORT" required:"true"`
	DB_USER         string        `envconfig:"DB_USER" required:"true"`
	DB_PASSWORD     string        `envconfig:"DB_PASSWORD" required:"true"`
	DB_NAME         string        `envconfig:"DB_NAME" required:"true"`
	JWT_SECRET      string        `envconfig:"JWT_SECRET" required:"true"`
	JWT_ACCESS_TTL  time.Duration `envconfig:"JWT_ACCESS_TTL" required:"true"`
	JWT_REFRESH_TTL time.Duration `envconfig:"JWT_REFRESH_TTL" required:"true"`
}

func main() {
	app := fiber.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app.Use(loggerMiddleware.Logger(logger))

	var env Env
	if err := envconfig.Process("", &env); err != nil {
		slog.Error("failed to process variables from env:", "error", err)
		return
	}

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", env.DB_USER, env.DB_PASSWORD, env.DB_HOST, env.DB_PORT, env.DB_NAME))
	if err != nil {
		slog.Error("postgres connect failed", "error", err)
		return
	}
	defer pool.Close()
	queries := sqlc.New(pool)

	authRepo := auth.NewRepository(queries)
	userRepo := user.NewRepository(queries)

	authService := auth.NewService(authRepo, userRepo, env.JWT_ACCESS_TTL, env.JWT_REFRESH_TTL, env.JWT_SECRET)
	userService := user.NewService(userRepo)

	authMiddleware := auth.AuthMiddleware(authService)

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)

	authHandler.Register(app)
	userHandler.Register(app, authMiddleware)

	logger.Info("starting API")
	log.Fatal(app.Listen(":" + env.APP_PORT))
}
