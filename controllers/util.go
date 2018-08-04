package controllers

import (
	"net/http"
)

// Write common http error
func WriteHttpError(err error, w http.ResponseWriter) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
