package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"math/rand"

	"github.com/gorilla/mux"
	
	"url_shortener/db"
	"url_shortener/models"
)

var (
	collectionName = "shorturls"
)

func generateShortKey() string {
    characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    keyLength := 6
    var result strings.Builder
    for i := 0; i < keyLength; i++ {
        result.WriteByte(characters[rand.Intn(len(characters))])
    }
    return result.String()
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	originalURL := r.FormValue("url")

	shortKey := generateShortKey()

	collection := db.MongoClient.Database(os.Getenv("DATABASE_NAME")).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), models.URL{
		ShortKey:    shortKey,
		OriginalURL: originalURL,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	shortenedURL := os.Getenv("BASE_URL") + shortKey
	fmt.Fprintf(w, "Shortened URL: %s", shortenedURL)
}

func RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortKey := params["shortURL"]

	collection := db.MongoClient.Database(os.Getenv("DATABASE_NAME")).Collection(collectionName)
	var result models.URL
	err := collection.FindOne(context.Background(), map[string]interface{}{
		"shortKey": shortKey,
	}).Decode(&result)

	if err != nil {
		log.Println(err)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, result.OriginalURL, http.StatusFound)
}
