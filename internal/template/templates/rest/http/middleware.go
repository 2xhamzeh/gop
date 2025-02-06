package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"example.com/rest"
	"github.com/google/uuid"
)

// Custom type for context keys to prevent collisions
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	requestIDKey contextKey = "request_id"
)

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(validate func(string) (int, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractBearerToken(r)
			if err != nil {
				writeError(w, err)
				return
			}

			userID, err := validate(token)
			if err != nil {
				writeError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractBearerToken handles Authorization header parsing and validation
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", rest.Errorf(rest.UNAUTHORIZED_ERROR, "authorization header required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", rest.Errorf(rest.UNAUTHORIZED_ERROR, "invalid authorization header format")
	}

	return parts[1], nil
}

// GetUserID safely retrieves user ID from context with proper type assertion
func getUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0, rest.Errorf(rest.UNAUTHORIZED_ERROR, "user not authenticated")
	}
	return userID, nil
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID, _ := r.Context().Value(requestIDKey).(string)

		slog.Info("HTTP Request Received",
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(w, r)

		slog.Info("HTTP Request Completed",
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
			"user_agent", r.UserAgent(),
		)
	})
}

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				slog.Error("panic occurred", "error", err, "request_id", r.Context().Value(requestIDKey))
				writeError(w, rest.Errorf(rest.INTERNAL_ERROR, "internal server error"))
			}
		}()

		next.ServeHTTP(w, r)

	})
}

// Not found handler
func notFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, rest.Errorf(rest.NOTFOUND_ERROR, "resource not found"))
}

// Method not allowed handler
func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeError(w, rest.Errorf(rest.INVALID_ERROR, "method not allowed"))
}
