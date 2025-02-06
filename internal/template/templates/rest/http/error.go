package http

import (
	"encoding/json"
	"example.com/rest"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Message string   `json:"message"`
	Fields  []string `json:"fields,omitempty"`
}

var codes = map[string]int{
	rest.CONFLICT_ERROR:     http.StatusConflict,
	rest.INVALID_ERROR:      http.StatusBadRequest,
	rest.NOTFOUND_ERROR:     http.StatusNotFound,
	rest.UNAUTHORIZED_ERROR: http.StatusUnauthorized,
	rest.INTERNAL_ERROR:     http.StatusInternalServerError,
}

func writeError(w http.ResponseWriter, err error) {
	errResp := errorResponse{
		Message: rest.ErrorMessage(err),
		Fields:  rest.ErrorFields(err),
	}
	status := codes[rest.ErrorCode(err)]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]errorResponse{"error": errResp})

	slog.Info("responded with error", "error", errResp.Message, "status", status)
}
