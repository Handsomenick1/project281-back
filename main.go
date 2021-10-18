package main

import (
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
    jwt "github.com/form3tech-oss/jwt-go"

)

func main() {
	fmt.Println("started-service")
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            return []byte(mySigningKey), nil
        },
        SigningMethod: jwt.SigningMethodHS256,
    })

	r := mux.NewRouter()
	r.Handle("/upload",  jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST", "OPTIONS")
	r.Handle("/delete", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("POST", "OPTIONS")
	r.Handle("/get", 	jwtMiddleware.Handler(http.HandlerFunc(getHandler))).Methods("GET", "OPTIONS")
	r.Handle("/update", jwtMiddleware.Handler(http.HandlerFunc(updateHandler))).Methods("POST", "OPTIONS")
	r.Handle("/signup", http.Handler(http.HandlerFunc(signupHandler))).Methods("POST", "OPTIONS")
    r.Handle("/signin", http.Handler(http.HandlerFunc(signinHandler))).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", r))
}

