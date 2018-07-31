package controllers

import (
	"encoding/json"
	"github.com/dtsang7/ASAPP/models"
	"net/http"
)

type Auth struct {
	Id    int
	Token string
}

//handles creating new user
func (h Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	var usr models.User
	json.NewDecoder(r.Body).Decode(&usr)

	//validate
	err := ValidateUser(usr)
	if err != nil {
		WriteJsonError(err, w)
		return
	}

	id, err := h.DB.CreateUser(usr)

	if err != nil {
		WriteJsonError(err, w)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]int{"id": id})
	if err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}

//handles login of existing user
func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	//get info from body and put in User struct
	var existingUser models.User
	json.NewDecoder(r.Body).Decode(&existingUser)

	//validate
	err := ValidateUser(existingUser)
	if err != nil {
		WriteJsonError(err, w)
		return
	}

	//authenticate user
	id, tokenString, err := h.Authenticate(existingUser)
	if err != nil {
		WriteJsonError(err, w)
		return
	}

	auth := Auth{id, tokenString}
	err = json.NewEncoder(w).Encode(auth)
	if err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}
