package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DavAnders/odin-blogapi/backend/internal/model"
	"github.com/DavAnders/odin-blogapi/backend/internal/repository"
	"github.com/DavAnders/odin-blogapi/backend/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

type UserController struct {
	repo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{
		repo: repo,
	}
}

// Create user with JWT token
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
    var user model.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := c.repo.CreateUser(r.Context(), user); err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    token, err := jwt.GenerateToken(user)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // setTokenAsCookie(w, token)
	// Return the token in the JSON response instead of setting a cookie for now
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}


// Handles POST requests to create a new user
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user model.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, fmt.Sprintf("Invalid user data: %v", err), http.StatusBadRequest)
        return
    }

    if err := c.repo.CreateUser(context.Background(), user); err != nil {
        log.Printf("Failed to create user: %v", err)  
        http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// Handles GET requests to retrieve a single user
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := c.repo.GetUser(context.Background(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Handles GET requests to retrieve all users
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.repo.GetUsers(context.Background())
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}