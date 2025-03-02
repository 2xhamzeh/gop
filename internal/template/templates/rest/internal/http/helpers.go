package http

import (
	"net/http"

	"example.com/rest/internal/domain"
)

// getUserID safely retrieves user ID from context
func getUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0, domain.Errorf(domain.UNAUTHORIZED_ERROR, "user not authenticated")
	}
	return userID, nil
}
