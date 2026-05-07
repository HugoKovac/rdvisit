package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/HugoKovac/rdvisit/db/sqlc"
	"github.com/HugoKovac/rdvisit/internal/appointment"
	"github.com/HugoKovac/rdvisit/internal/auth"
	"github.com/HugoKovac/rdvisit/internal/center"
	"github.com/HugoKovac/rdvisit/internal/user"
	loggerMiddleware "github.com/HugoKovac/rdvisit/pkg/logger/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type Env struct {
	APP_PORT              string        `envconfig:"APP_PORT" required:"true"`
	DB_HOST               string        `envconfig:"DB_HOST" required:"true"`
	DB_PORT               string        `envconfig:"DB_PORT" required:"true"`
	DB_USER               string        `envconfig:"DB_USER" required:"true"`
	DB_PASSWORD           string        `envconfig:"DB_PASSWORD" required:"true"`
	DB_NAME               string        `envconfig:"DB_NAME" required:"true"`
	JWT_SECRET            string        `envconfig:"JWT_SECRET" required:"true"`
	JWT_ACCESS_TTL        time.Duration `envconfig:"JWT_ACCESS_TTL" required:"true"`
	JWT_REFRESH_TTL       time.Duration `envconfig:"JWT_REFRESH_TTL" required:"true"`
	MJ_APIKEY_PUBLIC      string        `envconfig:"MJ_APIKEY_PUBLIC" required:"true"`
	MJ_APIKEY_PRIVATE     string        `envconfig:"MJ_APIKEY_PRIVATE" required:"true"`
	MJ_SENDER_EMAIL       string        `envconfig:"MJ_SENDER_EMAIL" required:"true"`
	MJ_APP_NAME           string        `envconfig:"MJ_APP_NAME" required:"true"`
	VERIFICATION_CODE_TTL time.Duration `envconfig:"VERIFICATION_CODE_TTL" required:"true"`
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

	client := mailjet.NewMailjetClient(env.MJ_APIKEY_PUBLIC, env.MJ_APIKEY_PRIVATE)

	userMailer := user.NewMailer(client, env.MJ_APP_NAME, env.MJ_SENDER_EMAIL)

	authRepo := auth.NewRepository(queries)
	userRepo := user.NewRepository(queries)
	centerRepo := center.NewRepository(queries)
	appointmentRepo := appointment.NewRepository(queries)

	authService := auth.NewService(authRepo, userRepo, userMailer, env.JWT_ACCESS_TTL, env.JWT_REFRESH_TTL, env.VERIFICATION_CODE_TTL, env.JWT_SECRET)
	userService := user.NewService(userRepo, userMailer)
	centerService := center.NewService(centerRepo)
	appointmentService := appointment.NewService(appointmentRepo)

	authMiddleware := auth.AuthMiddleware(authService)
	centerMiddleware := center.CheckNGetCenter(centerService)
	verifiedMiddleware := user.CheckUserVerified(userService)

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	centerHandler := center.NewHandler(centerService)
	appointmentHandler := appointment.NewHandler(appointmentService)

	authHandler.Register(app)
	userHandler.Register(app, authMiddleware)
	centerHandler.Register(app, authMiddleware, centerMiddleware)
	appointmentHandler.Register(app, authMiddleware, verifiedMiddleware)

	logger.Info("starting API")
	log.Fatal(app.Listen(":" + env.APP_PORT))
}
