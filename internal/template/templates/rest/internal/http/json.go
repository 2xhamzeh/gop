package http

import (
	"encoding/json"
	"net/http"

	"example.com/rest/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func decodeJSON(r *http.Request, target any) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return domain.Errorf(domain.INVALID_ERROR, "invalid request body")
	}
	return nil
}
