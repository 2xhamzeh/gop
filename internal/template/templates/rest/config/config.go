// config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	Environment string
	DB          dBConfig
	Server      serverConfig
	JWT         jWTConfig
}

type dBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type serverConfig struct {
	Host string
	Port int
}

type jWTConfig struct {
	Secret   string
	Duration time.Duration
}

func Load() (*config, error) {
	godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	jwtDuration, err := time.ParseDuration(getEnv("JWT_DURATION", "1h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_DURATION: %w", err)
	}

	return &config{
		Environment: getEnv("APP_ENV", "development"),
		DB: dBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "2xh"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "notesapp"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: serverConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: serverPort,
		},
		JWT: jWTConfig{
			Secret:   getEnv("JWT_SECRET", ""),
			Duration: jwtDuration,
		},
	}, nil
}

// GetDSN returns the formatted database connection string
func (c *dBConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

// GetAddress returns the formatted server address
func (c *serverConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// getEnv retrieves an environment variable value or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
