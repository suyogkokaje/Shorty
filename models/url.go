package models

type URL struct {
	ShortKey    string `bson:"shortKey"`
	OriginalURL string `bson:"originalURL"`
}
