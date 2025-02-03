package http

import (
	"encoding/json"
	"example.com/app"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Message string   `json:"message"`
	Fields  []string `json:"fields,omitempty"`
}

var codes = map[string]int{
	app.CONFLICT_ERROR:     http.StatusConflict,
	app.INVALID_ERROR:      http.StatusBadRequest,
	app.NOTFOUND_ERROR:     http.StatusNotFound,
	app.UNAUTHORIZED_ERROR: http.StatusUnauthorized,
	app.INTERNAL_ERROR:     http.StatusInternalServerError,
}

func writeError(w http.ResponseWriter, err error) {
	errResp := errorResponse{
		Message: app.ErrorMessage(err),
		Fields:  app.ErrorFields(err),
	}
	status := codes[app.ErrorCode(err)]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]errorResponse{"error": errResp})

	slog.Info("responded with error", "error", errResp.Message, "status", status)
}
