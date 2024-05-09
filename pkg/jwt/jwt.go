package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// Generates a new JWT token
func GenerateToken(username string) (string, error) {
    expirationTime := time.Now().Add(1 * time.Hour) // Token expires after 1 hour
    claims := &Claims{
        Username: username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)

    return tokenString, err
}

// Validates a JWT token
func ValidateToken(tokenString string) (*jwt.Token, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    return token, err
}
