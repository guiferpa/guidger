package httputil

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, http.StatusInternalServerError, err)
	}
}

func Error(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
