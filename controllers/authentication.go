package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dtsang7/ASAPP/models"
	"time"
)

func (h Handler) Authenticate(user models.User) (int, string, error) {

	var tokenString string
	id, err := h.DB.LoginUser(user)
	if err != nil {
		return 0, tokenString, err
	}
	tokenString, err = createToken(id)
	if err != nil {
		return 0, tokenString, err
	}
	return id, tokenString, nil
}

func createToken(id int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 300).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
