package middlewares

import (
	"context"
	// "fmt"
	"net/http"

	"url_shortener/utils"
	// "github.com/gorilla/context"
)

// Authentication validates token and authorizes users
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			http.Error(w, "No Authorization header provided", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateToken(clientToken)
		if err != "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store user claims in the request context
		// context.Set(r, "userClaims", claims)

		ctx := context.WithValue(r.Context(), "userClaims", claims)
		r = r.WithContext(ctx)

		// fmt.Println(r.Context().Value("userClaims"))

		// Continue processing the request
		next.ServeHTTP(w, r)
	})
}
