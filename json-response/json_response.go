// Package jsonresponse has functions for returning JSONs in APIs
package jsonresponse

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON adds JSON response for API request
func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		panic(err)
	}
}

// RespondWithError adds JSON error response for API request
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
