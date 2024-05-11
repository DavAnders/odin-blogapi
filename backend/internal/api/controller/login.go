package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavAnders/odin-blogapi/backend/pkg/jwt"
)

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate user credentials (you need to implement this in your UserRepository)
    user, err := c.repo.ValidateCredentials(r.Context(), credentials.Username, credentials.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Generate JWT
    token, err := jwt.GenerateToken(*user)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Set security headers
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "no-store")
    w.Header().Set("Pragma", "no-cache")

    // Send token in HTTP-only cookie
    http.SetCookie(w, &http.Cookie{
        Name: "token",
        Value: token,
        Expires: time.Now().Add(1 * time.Hour),
        HttpOnly: true,
        Path: "/",
        Secure: false, // For development purposes only
        SameSite: http.SameSiteStrictMode,
    })

    // Send token in response body (alternative method)
    // json.NewEncoder(w).Encode(map[string]string{"token": token})
}
