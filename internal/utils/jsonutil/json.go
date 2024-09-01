package jsonutil

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func WriteErrorJSON(w http.ResponseWriter, statusCode int, message string) error {
	errResponse := map[string]string{"error": message}
	return WriteJSON(w, statusCode, errResponse)
}
