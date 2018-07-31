package main

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/dtsang7/ASAPP/controllers"
	"github.com/dtsang7/ASAPP/models"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

func main() {
	dao := models.CreateDAO("sqlite3", "challenge.db")

	dao.RunMigrations()

	handler := controllers.Handler{DB: dao}

	publicRouter := mux.NewRouter()
	protectedRouter := mux.NewRouter()

	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	//public
	publicRouter.HandleFunc("/check", handler.CheckHandler).Methods("POST")
	publicRouter.HandleFunc("/users", handler.UserHandler).Methods("POST")
	publicRouter.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	//protected
	protectedRouter.HandleFunc("/messages", handler.SendMessageHandler).Methods("POST")
	protectedRouter.HandleFunc("/messages", handler.GetMessagesHandler).Methods("GET")

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(protectedRouter))
	publicRouter.PathPrefix("/").Handler(an)

	n := negroni.Classic()
	n.UseHandler(publicRouter)
	log.Println("Starting server on 8080")
	http.ListenAndServe(":8080", n)
}
