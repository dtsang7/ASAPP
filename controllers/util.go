package controllers

import (
	"net/http"
)

func WriteHttpError(err error, w http.ResponseWriter) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
