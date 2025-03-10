package http

import (
	"log/slog"
	"net/http"

	"example.com/rest/internal/domain"
)

// baseHandler contains common dependencies for all handlers.
type baseHandler struct {
	json *jsonHelper
}

// NewBaseHandler creates a new base handler which contains common dependencies for all handlers.
func NewBaseHandler(logger *slog.Logger) *baseHandler {
	return &baseHandler{
		json: &jsonHelper{logger: logger},
	}
}

// getUserID safely retrieves user ID from the context
func (b *baseHandler) getUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0, domain.Errorf(domain.UNAUTHORIZED_ERROR, "user not authenticated")
	}
	return userID, nil
}
