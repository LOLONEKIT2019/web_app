package main

import (
	. "./handlers"
	"./config"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatal(errors.New("Config not found"))
	}

	r := mux.NewRouter()

	r.Handle("/status", http.HandlerFunc(RootHandler))

	r.Handle("/token", NewJwtToken(cfg))

	jwt := NewJwtAuthentication(cfg)

	r.Handle("/users", jwt.JwtAuthentication(NewUserRouter(cfg)))
	r.Handle("/users/{id}", jwt.JwtAuthentication(NewUserRouter(cfg)))
	r.Handle("/contacts/{id}", jwt.JwtAuthentication(http.HandlerFunc(ContactRouter)))

	log.Println(fmt.Sprintf("Api start on port : %v",cfg.Port))

	err := http.ListenAndServe(fmt.Sprintf(":%v",cfg.Port), r)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
