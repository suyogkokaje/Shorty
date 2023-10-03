package utils

import (
	"github.com/google/uuid"
	"math/big"
)

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func base62Encode(s string) string {
	base := len(base62Chars)
	num, _ := new(big.Int).SetString(s, 16)
	encoded := make([]byte, 0)
	for num.Cmp(big.NewInt(0)) > 0 {
		quotient, remainder := new(big.Int).DivMod(num, big.NewInt(int64(base)), new(big.Int))
		encoded = append(encoded, base62Chars[remainder.Int64()])
		num = quotient
	}
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}
	return string(encoded)
}

func GenerateShortKey() string {
	uuid := uuid.New()
	shortKey := uuid.String()[:6]
	shortKey = base62Encode(shortKey)
	return shortKey
}
