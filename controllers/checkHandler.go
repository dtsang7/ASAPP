package controllers

import (
	"encoding/json"
	"github.com/dtsang7/ASAPP/models"
	"net/http"
)

type Handler struct {
	DB *models.DAO
}

//handles check
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

	err = json.NewEncoder(w).Encode(map[string]string{"health": "ok"})
	if err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}
