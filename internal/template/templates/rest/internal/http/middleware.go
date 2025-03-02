package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"example.com/rest/internal/domain"
	"github.com/google/uuid"
)

type middlewares struct {
	validateToken func(string) (int, error)
	logger        *slog.Logger
}

func NewMiddlewares(validateToken func(string) (int, error), logger *slog.Logger) *middlewares {
	return &middlewares{
		validateToken: validateToken,
		logger:        logger,
	}
}

// Custom type for context keys to prevent collisions
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	requestIDKey contextKey = "request_id"
)

// NewAuthMiddleware creates a new authentication middleware.
// It takes in the token validation function and returns a middleware function.
func (m *middlewares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, domain.Errorf(domain.UNAUTHORIZED_ERROR, "authorization header required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeError(w, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid authorization header format"))
			return
		}

		userID, err := m.validateToken(parts[1])
		if err != nil {
			writeError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requestID middleware generates a unique request ID and adds it to the request context
func (m *middlewares) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// loggerMiddleware logs incoming HTTP requests and their duration.
// It logs twice: when the request is received and when it is completed.
func (m *middlewares) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		reqID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			reqID = "unknown"
		}

		m.logger.Info("HTTP Request Received",
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(wrapped, r)

		m.logger.Info("HTTP Request Completed",
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.status,
			"duration", time.Since(start).String(),
			"user_agent", r.UserAgent(),
		)
	})
}

// recovery middleware recovers from panics and logs the error.
func (m *middlewares) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.logger.Error("panic occurred", "error", err, "request_id", r.Context().Value(requestIDKey))
				writeError(w, domain.Errorf(domain.INTERNAL_ERROR, "internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
