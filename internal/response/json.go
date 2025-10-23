// Package response provides utilities for handling JSON HTTP responses and requests.
package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response with the given status code and data
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Error writes a JSON error response with the given status code and message
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, map[string]string{
		"error": message,
	})
}

// DecodeJSON decodes the JSON body from the request into the provided value
func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
