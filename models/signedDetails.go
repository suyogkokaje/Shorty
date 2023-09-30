package models

import jwt "github.com/dgrijalva/jwt-go"

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}