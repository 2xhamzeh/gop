package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"example.com/app"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err, "status", status, "data", data)
		writeError(w, app.Errorf(app.INTERNAL_ERROR, "internal server error"))
		return
	}
	slog.Info("responded with success", "status", status, "data", data)
}

func decodeJSON(r *http.Request, target any) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return app.Errorf(app.INVALID_ERROR, "invalid request body")
	}
	return nil
}
