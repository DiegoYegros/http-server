package handlers

import (
	"encoding/json"
	"net/http"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	response := map[string]interface{}{
		"message": "This is my website.",
		"params":  query,
	}
	sendJSONResponse(w, response, http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
