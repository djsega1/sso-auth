package main

import (
	"log"
	"net/http"

	"github.com/djsega1/sso-auth/config"
	"github.com/djsega1/sso-auth/database"
	"github.com/djsega1/sso-auth/handler"
	"github.com/djsega1/sso-auth/repository"
	"github.com/djsega1/sso-auth/service"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	database.InitDB(cfg)

	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService, cfg)

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/validate", authHandler.Validate).Methods("GET")

	log.Println("Auth REST API server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
