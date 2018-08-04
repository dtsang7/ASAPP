package main

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/dtsang7/ASAPP/config"
	"github.com/dtsang7/ASAPP/controllers"
	"github.com/dtsang7/ASAPP/models"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// load config
	env := strings.ToLower(os.Getenv("ASAPP_ENV"))
	config, err := config.GetConfig(env)
	if err != nil {
		log.Println("Server fail to start, unable to retrieve config")
		os.Exit(1)
	}
	// connect to data store
	dao := models.CreateDAO(config.DBDriver, config.DBName)
	dao.RunMigrations()

	// Set up router
	handler := controllers.Handler{DB: dao}
	publicRouter := mux.NewRouter()
	protectedRouter := mux.NewRouter()

	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// This should be an actual secret
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	//public
	publicRouter.HandleFunc("/check", handler.CheckHandler).Methods("POST")
	publicRouter.HandleFunc("/users", handler.UserHandler).Methods("POST")
	publicRouter.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	//protected (jwt)
	protectedRouter.HandleFunc("/messages", handler.SendMessageHandler).Methods("POST")
	protectedRouter.HandleFunc("/messages", handler.GetMessagesHandler).Methods("GET")

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(protectedRouter))
	publicRouter.PathPrefix("/").Handler(an)

	n := negroni.Classic()
	n.UseHandler(publicRouter)
	log.Println("Starting server on " + config.Port)
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, n))
}
