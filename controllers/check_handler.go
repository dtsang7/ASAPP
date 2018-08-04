package controllers

import (
	"encoding/json"
	"github.com/dtsang7/ASAPP/models"
	"net/http"
)

type Health struct {
	Health string `json:"health"`
}

type Handler struct {
	DB        *models.DAO
	JWTSecret string
}

// checks system health
func (h Handler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	var res int
	res, err := h.DB.CheckDB()

	if err != nil {
		http.Error(w, "DB connection error", http.StatusInternalServerError)
		return
	}

	if res != 1 {
		http.Error(w, "Unexpected query result", http.StatusInternalServerError)
		return
	}

	health := Health{"ok"}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(health)
	if err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}
