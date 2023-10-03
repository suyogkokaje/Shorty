package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"url_shortener/db"
	"url_shortener/models"
	"url_shortener/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName = "shorturls"
)

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	originalURL := r.FormValue("url")
	password := r.FormValue("password")
	customShortKey := r.FormValue("customShortKey")

	userClaims, ok := r.Context().Value("userClaims").(*models.SignedDetails)
	if !ok {
		http.Error(w, "Failed to get user claims", http.StatusInternalServerError)
		return
	}

	collection := db.MongoClient.Database(os.Getenv("DATABASE_NAME")).Collection(collectionName)

	if customShortKey != "" {
		count, err := collection.CountDocuments(context.Background(), bson.M{"shortKey": customShortKey})
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to check custom short key uniqueness", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "Custom short key is already in use", http.StatusBadRequest)
			return
		}
	} else {
		customShortKey = utils.GenerateShortKey()
	}

	_, err := collection.InsertOne(context.Background(), models.URL{
		ShortKey:         customShortKey,
		OriginalURL:      originalURL,
		UserID:           userClaims.Uid,
		ClickCount:       0,
		CreatedAt:        time.Now(),
		LastRedirectedAt: time.Time{},
		Password:         password,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	shortenedURL := os.Getenv("BASE_URL") + customShortKey
	fmt.Fprintf(w, "Shortened URL: %s", shortenedURL)
}

func RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortKey := params["shortURL"]

	password := r.FormValue("password")

	collection := db.MongoClient.Database(os.Getenv("DATABASE_NAME")).Collection(collectionName)
	var result models.URL
	err := collection.FindOneAndUpdate(context.Background(), bson.M{"shortKey": shortKey},
		bson.M{"$inc": bson.M{"clickCount": 1}, "$set": bson.M{"lastRedirectedAt": time.Now()}}).Decode(&result)

	if err != nil {
		log.Println(err)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	if result.Password != "" {
		if password != result.Password {
			http.Error(w, "Access denied. Invalid password.", http.StatusUnauthorized)
			return
		}
	}

	redirectReq, _ := http.NewRequest("GET", result.OriginalURL, nil)
	redirectReq.Header = r.Header

	http.Redirect(w, r, redirectReq.URL.String(), http.StatusFound)
}
