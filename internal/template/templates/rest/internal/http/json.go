package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"example.com/rest/internal/domain"
	"example.com/rest/internal/validator"
)

// jsonHelper for encoding and decoding JSON.
type jsonHelper struct {
	logger *slog.Logger
}

// response is the JSON response format for the API.
type response struct {
	Status  string `json:"status"`            // "success" or "error"
	Message string `json:"message,omitempty"` // error message
	Data    any    `json:"data,omitempty"`    // response data or error fields for invalid requests
}

// WriteResponse sends a JSON response with the status code and response.
func (j *jsonHelper) WriteResponse(w http.ResponseWriter, statusCode int, res response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		// This error is very rare and should not happen
		// If it does, it means that the encoded data is not supported by the JSON encoder (e.g. chan or func)
		j.logger.Error("failed to encode response", "error", err)
	}
}

// Write sends a JSON response with the data and status code.
// It should be used in the case of success.
func (j *jsonHelper) Write(w http.ResponseWriter, statusCode int, data any) {
	response := response{
		Status: "success",
		Data:   data,
	}
	j.WriteResponse(w, statusCode, response)
}

// WriteError sends a JSON response with the error.
// If the error is a domain error, it sends the error message and status code.
// If the error is a validation error, it sends the error message, status code and fields.
// If the error is neither, it sends a generic error message and status code (500).
// Internal server errors are logged.
func (j *jsonHelper) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	res := response{Status: "error"}
	statusCode := http.StatusInternalServerError

	// Check if the error is a domain or validation error
	var domainError *domain.Error
	var validationError *validator.Error
	if errors.As(err, &domainError) {
		if code, ok := domainToHTTPErrors[domainError.Code]; ok {
			statusCode = code
			res.Message = domainError.Message
			res.Data = domainError.Fields
		}
	} else if errors.As(err, &validationError) {
		statusCode = http.StatusBadRequest
		res.Message = validationError.Message
		res.Data = validationError.Fields
	}

	if statusCode == http.StatusInternalServerError {
		// Log the error, this is the only place where we log errors
		reqID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			reqID = "unknown"
		}
		j.logger.Error("internal server error", "request_id", reqID, "error", err)
		// Hide the error message from the user
		res.Message = "internal server error"
	}
	j.WriteResponse(w, statusCode, res)
}

// Read decodes the JSON request body into the target.
// It doesn't allow unknown fields and returns an error if the request body is invalid.
func (j *jsonHelper) Read(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return domain.Errorf(domain.INVALID_ERROR, "invalid request body")
	}
	return nil
}
