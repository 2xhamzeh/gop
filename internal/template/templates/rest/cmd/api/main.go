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
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("server failed", "error", err)
	}
}

func run(logger *slog.Logger) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Connect to database
	db, err := postgres.New(cfg.DB.URL)
	if err != nil {
		return err
	}

	// Initialize services
	userService := postgres.NewUserService(db)
	authService := jwt.NewAuthService(cfg.JWT.Secret, cfg.JWT.Duration)

	// Initialize handlers and middlewares
	baseHandler := http.NewBaseHandler(logger)
	userHandler := http.NewUserHandler(baseHandler, userService, authService.GenerateToken)
	middlewares := http.NewMiddlewares(baseHandler, authService.ValidateToken, logger)

	// Initialize router
	router := http.NewRouter(userHandler, middlewares)

	// Start server
	server := http.NewServer(cfg.Server.Addr(), router, logger)
	return server.Start()
}
