package render

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func writeErrorJSON(w http.ResponseWriter, statusCode int, err string) error {
	type Envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, statusCode, &Envelope{Error: err})
}
