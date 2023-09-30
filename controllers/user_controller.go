package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"url_shortener/db"
	"url_shortener/models"
	"url_shortener/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = db.OpenCollection(db.DBinstance(), "user")
var validate = validator.New()

func HashPassword(password string) string {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        log.Panic(err)
    }

    return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
    err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
    check := true
    msg := ""

    if err != nil {
        msg = fmt.Sprintf("login or passowrd is incorrect")
        check = false
    }

    return check, msg
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    validationErr := validate.Struct(user)
    if validationErr != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": validationErr.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()

    count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
    if err != nil {
        log.Panic(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "error occurred while checking for the email"})
        return
    }

    password := HashPassword(user.Password)
    user.Password = password

    count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
    if err != nil {
        log.Panic(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "error occurred while checking for the phone number"})
        return
    }

    if count > 0 {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "this email or phone number already exists"})
        return
    }

    user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    user.ID = primitive.NewObjectID()
    user.User_id = user.ID.Hex()
    token, refreshToken, _ := utils.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.User_id)
    user.Token = token
    user.Refresh_token = refreshToken

    resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
    if insertErr != nil {
        msg := fmt.Sprintf("User item was not created")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": msg})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(resultInsertionNumber)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    var foundUser models.User

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()

    err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "login or password is incorrect"})
        return
    }

    passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
    if !passwordIsValid {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": msg})
        return
    }

    token, refreshToken, _ := utils.GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_id)

    utils.UpdateAllTokens(token, refreshToken, foundUser.User_id)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(foundUser)
}

