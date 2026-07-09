package util

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes a JSON response with the given status code and data.
func WriteJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
