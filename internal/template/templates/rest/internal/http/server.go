package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

// server represents an HTTP server.
type server struct {
	server *http.Server
	logger *slog.Logger
}

// NewServer creates a new HTTP server. It takes in an address in the format "host:port",
// a router, and a logger.
func NewServer(addr string, router http.Handler, logger *slog.Logger) *server {
	return &server{
		server: &http.Server{
			Addr:         addr,
			Handler:      router,
			IdleTimeout:  defaultIdleTimeout,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		logger: logger,
	}
}

// Start starts the HTTP server. It listens for incoming requests and blocks until the server is stopped.
// It also listens for the interrupt signal and gracefully shuts down the server.
// It returns an error if the server fails to start or shutdown.
func (s *server) Start() error {
	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.Error("failed to shutdown server", "error", err)
		}
		close(done)
	}()

	s.logger.Info("starting server", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	<-done
	s.logger.Info("server stopped gracefully")
	return nil
}
