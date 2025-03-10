package http

// Custom type for context keys
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	requestIDKey contextKey = "request_id"
)
