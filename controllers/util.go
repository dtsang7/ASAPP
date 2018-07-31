package controllers

import (
	"encoding/json"
	"net/http"
)

func WriteJsonError(err error, w http.ResponseWriter) {
	jsonErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if jsonErr != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}
