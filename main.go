package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"url_shortener/controllers"
	"url_shortener/db"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitMongoClient()

	r := mux.NewRouter()
	r.HandleFunc("/shorten", controllers.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", controllers.RedirectToOriginal).Methods("GET")

	http.Handle("/", r)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
