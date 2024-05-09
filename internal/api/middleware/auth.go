package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Auth middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate alg
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid Token: %v", err), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, fmt.Sprintf("Invalid Token: %v", err), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}))
}