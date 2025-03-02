package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"example.com/rest/internal/domain"
)

type errorResponse struct {
	Message string   `json:"message"`
	Fields  []string `json:"fields,omitempty"`
}

var codes = map[string]int{
	domain.CONFLICT_ERROR:     http.StatusConflict,
	domain.INVALID_ERROR:      http.StatusBadRequest,
	domain.NOTFOUND_ERROR:     http.StatusNotFound,
	domain.UNAUTHORIZED_ERROR: http.StatusUnauthorized,
	domain.INTERNAL_ERROR:     http.StatusInternalServerError,
}

// writeError sends a JSON response with the error message and fields.
// It takes in a domain error and writes the error message
// and fields to the response with the appropriate status code.
func writeError(w http.ResponseWriter, err error) {
	res := errorResponse{}
	var status int
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		res.Message = domainErr.Message
		res.Fields = domainErr.Fields
		status = codes[domainErr.Code]
	} else {
		res.Message = "internal server error"
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]errorResponse{"error": res})
}

// writeHTTPError sends a JSON response with the error message and status code.
func writeHTTPError(w http.ResponseWriter, status int, message string) {
	res := errorResponse{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]errorResponse{"error": res})
}
