package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"url_shortener/controllers"
	"url_shortener/db"
	"url_shortener/middlewares"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitMongoClient()

	r := mux.NewRouter()

	r.HandleFunc("/signup", controllers.SignUpHandler).Methods("POST")
	r.HandleFunc("/login", controllers.LoginHandler).Methods("POST")

	shortenerRoute := r.PathPrefix("/shorten").Subrouter()
	shortenerRoute.Use(middlewares.Authentication)
	shortenerRoute.HandleFunc("", controllers.ShortenURL).Methods("POST")

	redirectRoute := r.PathPrefix("/{shortURL}").Subrouter()
	redirectRoute.HandleFunc("", controllers.RedirectToOriginal).Methods("GET")

	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server started on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
