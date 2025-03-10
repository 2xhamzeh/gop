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
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	// Load environment variables
	godotenv.Load()
	// Load database URL
	DB_URL, err := config.LoadDB_URL()
	if err != nil {
		return err
	}
	// Create migrator
	m, err := migrate.New("file://migrations", DB_URL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	// Parse command line arguments and run the appropriate command
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("command required: up, down, or version")
	}
	command := args[0]

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
