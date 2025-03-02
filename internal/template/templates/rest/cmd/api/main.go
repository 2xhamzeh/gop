package main

import (
	"log/slog"
	"os"

	"example.com/rest/internal/config"
	"example.com/rest/internal/http"
	"example.com/rest/internal/jwt"
	"example.com/rest/internal/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	if err := run(logger); err != nil {
		logger.Error("server failed", "error", err)
	}
}

func run(logger *slog.Logger) error {

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger.Info("loaded configuration")

	db, err := postgres.New(cfg.DB.GetDSN())
	if err != nil {
		return err
	}

	logger.Info("connected to database")

	userService := postgres.NewUserService(db, logger)
	authService := jwt.NewAuthService(cfg.JWT.Secret, cfg.JWT.Duration, logger)
	noteService := postgres.NewNoteService(db, logger)

	logger.Info("initialized services")

	userHandler := http.NewUserHandler(userService, authService.GenerateToken)
	middlewares := http.NewMiddlewares(authService.ValidateToken, logger)
	noteHandler := http.NewNoteHandler(noteService)

	logger.Info("initialized handlers")

	router := http.NewRouter(userHandler, noteHandler, middlewares)

	logger.Info("initialized routes")

	server := http.New(cfg.Server.GetAddress(), router, logger)
	return server.Start()
}
