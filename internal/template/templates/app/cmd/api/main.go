package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"example.com/app/config"
	"example.com/app/http"
	"example.com/app/jwt"
	"example.com/app/postgres"
	"github.com/lmittmann/tint"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
}

func run() error {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			AddSource:  true,
		}),
	))

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	slog.Info("loaded configuration")

	db, err := postgres.New(cfg.DB.GetDSN())
	if err != nil {
		return err
	}

	slog.Info("connected to database")

	userService := postgres.NewUserService(db)
	noteService := postgres.NewNoteService(db)
	authService := jwt.NewAuthService(cfg.JWT.Secret, cfg.JWT.Duration)

	slog.Info("initialized services")

	userHandler := http.NewUserHandler(userService, authService.GenerateToken)
	noteHandler := http.NewNoteHandler(noteService)
	authMiddleware := http.NewAuthMiddleware(authService.ValidateToken)

	slog.Info("initialized handlers")

	router := http.NewRouter(userHandler, noteHandler, authMiddleware)

	slog.Info("initialized router")

	server := http.New(cfg.Server.GetAddress(), router)
	return server.Start()
}
