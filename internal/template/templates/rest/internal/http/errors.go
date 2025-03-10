package http

import (
	"net/http"

	"example.com/rest/internal/domain"
)

// domainToHTTPErrors maps domain error codes to HTTP status codes.
var domainToHTTPErrors = map[domain.ErrorCode]int{
	domain.INVALID_ERROR:      http.StatusBadRequest,          // 400
	domain.UNAUTHORIZED_ERROR: http.StatusUnauthorized,        // 401
	domain.FORBIDDEN_ERROR:    http.StatusForbidden,           // 403
	domain.NOTFOUND_ERROR:     http.StatusNotFound,            // 404
	domain.CONFLICT_ERROR:     http.StatusConflict,            // 409
	domain.INTERNAL_ERROR:     http.StatusInternalServerError, // 500
}
