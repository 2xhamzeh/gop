// config/config.go
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// config represents the application configuration
type config struct {
	DB     db
	Server server
	JWT    jwt
}

type db struct {
	URL string
}

type server struct {
	Host string
	Port string
}

type jwt struct {
	Secret   string
	Duration time.Duration
}

/*
Load loads the configuration from environment variables.

It returns a configuration struct or an error if a required variable is missing.

It expects the following environment variables:

	PGUSER
	PGPASSWORD
	PGHOST
	PGDATABASE
	PGSSLMODE

	JWT_SECRET
	JWT_DURATION

	SERVER_HOST
	SERVER_PORT
*/
func Load() (*config, error) {
	// Load environment variables, can be omitted if you don't use .env file and inject all variables via environment
	godotenv.Load()

	// Load database URL
	DB_URL, err := LoadDB_URL()
	if err != nil {
		return nil, err
	}

	// Load JWT configuration
	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	JWT_DURATION, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {
		return nil, fmt.Errorf("JWT_DURATION is invalid")
	}

	// Load server configuration
	SERVER_HOST := os.Getenv("SERVER_HOST")
	if SERVER_HOST == "" {
		return nil, fmt.Errorf("SERVER_HOST is required")
	}

	SERVER_PORT := os.Getenv("SERVER_PORT")
	if SERVER_PORT == "" {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}

	// Return configuration
	return &config{
		DB: db{
			URL: DB_URL,
		},
		Server: server{
			Host: SERVER_HOST,
			Port: SERVER_PORT,
		},
		JWT: jwt{
			Secret:   JWT_SECRET,
			Duration: JWT_DURATION,
		},
	}, nil
}

// Addr returns the server address in the format host:port
func (s *server) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

/*
LoadDB_URL loads the database URL from environment variables.

URL is returned in the format: postgres://user:password@host/database?sslmode=mode

It expects the following environment variables:

	PGUSER
	PGPASSWORD
	PGHOST
	PGDATABASE
	PGSSLMODE
*/
func LoadDB_URL() (string, error) {
	PGHOST := os.Getenv("PGHOST")

	PGUSER := os.Getenv("PGUSER")
	if PGUSER == "" {
		return "", fmt.Errorf("PGUSER is required")
	}

	PGPASSWORD := os.Getenv("PGPASSWORD")
	if PGPASSWORD == "" {
		return "", fmt.Errorf("PGPASSWORD is required")
	}

	if PGHOST == "" {
		return "", fmt.Errorf("PGHOST is required")
	}

	PGDATABASE := os.Getenv("PGDATABASE")
	if PGDATABASE == "" {
		return "", fmt.Errorf("PGDATABASE is required")
	}

	PGSSLMODE := os.Getenv("PGSSLMODE")
	if PGSSLMODE == "" {
		return "", fmt.Errorf("PGSSLMODE is required")
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", PGUSER, PGPASSWORD, PGHOST, PGDATABASE, PGSSLMODE), nil
}
