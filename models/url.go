package models

type URL struct {
	ShortKey    string `bson:"shortKey"`
	OriginalURL string `bson:"originalURL"`
	UserID      string `bson:"userID"`
}
