package models

import "time"

type URL struct {
	ShortKey         string    `bson:"shortKey"`
	OriginalURL      string    `bson:"originalURL"`
	UserID           string    `bson:"userID"`
	ClickCount       int       `bson:"clickCount" json:"clickCount"`
	CreatedAt        time.Time `bson:"createdAt" json:"createdAt"`
	LastRedirectedAt time.Time `bson:"lastRedirectedAt" json:"lastRedirectedAt"`
	Password         string    `json:"password" bson:"password"`
}
