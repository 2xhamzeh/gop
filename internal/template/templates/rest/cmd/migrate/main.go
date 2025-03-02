// cmd/migrate/main.go
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"example.com/rest/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	if err := run(logger); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {

	var dbURL string
	flag.StringVar(&dbURL, "db-url", "", "Database URL (overrides config)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("command required: up, down, or version")
	}
	command := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if dbURL == "" {
		dbURL = cfg.DB.GetDSN()
	}

	logger.Info("initializing migrator", "db_url", dbURL)

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run up migration: %w", err)
		}
		logger.Info("successfully ran up migrations")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run down migration: %w", err)
		}
		logger.Info("successfully ran down migrations")

	case "version":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return fmt.Errorf("failed to get version: %w", err)
		}
		logger.Info("current migration status",
			"version", version,
			"dirty", dirty,
		)

	default:
		return fmt.Errorf("invalid command %q: must be up, down, or version", command)
	}

	return nil
}
