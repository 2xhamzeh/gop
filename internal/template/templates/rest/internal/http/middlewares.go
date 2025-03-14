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

// Middlewares contains all the dependencies required by the middleware functions.
type Middlewares struct {
	*baseHandler
	validateToken func(string) (int, error)
	logger        *slog.Logger
}

// NewMiddlewares creates a new Middlewares instance with the required dependencies.
func NewMiddlewares(baseHandler *baseHandler, validateToken func(string) (int, error), logger *slog.Logger) *Middlewares {
	return &Middlewares{
		baseHandler:   baseHandler,
		validateToken: validateToken,
		logger:        logger,
	}
}

// Auth returns a middleware that validates the JWT token in the Authorization header.
// If the token is valid, it adds the user ID to the request context.
// Handlers can retrieve the user ID using the getUserID method from the baseHandler.
func (m *Middlewares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.json.WriteError(w, r, domain.Errorf(domain.UNAUTHORIZED_ERROR, "authorization header required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			m.json.WriteError(w, r, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid authorization header format"))
			return
		}

		userID, err := m.validateToken(parts[1])
		if err != nil {
			m.json.WriteError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestID middleware generates a unique request ID and adds it to the request context and response headers.
func (m *Middlewares) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader overrides the WriteHeader method to capture the status code.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger logs incoming HTTP requests and their duration.
// It logs at the start and end of a request.
func (m *Middlewares) Logger(next http.Handler) http.Handler {
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
			"user_agent", r.UserAgent(),
			"status", wrapped.status,
			"duration", time.Since(start).String(),
		)
	})
}

// recovery middleware recovers from panics and logs the error.
func (m *Middlewares) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.json.WriteError(w, r, domain.Errorf(domain.INTERNAL_ERROR, "panic occurred: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// NotFound sends a 404 response for unknown routes.
func (m *Middlewares) NotFound(w http.ResponseWriter, r *http.Request) {
	m.json.WriteResponse(w, http.StatusNotFound, response{
		Status:  "error",
		Message: "not found",
	})
}

// MethodNotAllowed sends a 405 response for unknown methods.
func (m *Middlewares) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	m.json.WriteResponse(w, http.StatusMethodNotAllowed, response{
		Status:  "error",
		Message: "method not allowed",
	})
}
