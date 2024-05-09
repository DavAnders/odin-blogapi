package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"github.com/DavAnders/odin-blogapi/internal/repository"
	"github.com/DavAnders/odin-blogapi/pkg/jwt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

    // Create user in the database
    if err := c.repo.CreateUser(r.Context(), user); err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    // Generate a JWT for the user after registration
    token, err := jwt.GenerateToken(user.Username)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}


// Handles POST requests to create a new user
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	user.ID = primitive.NewObjectID()  // Assuming MongoDB's ObjectIDs are used

	if err := c.repo.CreateUser(context.Background(), user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Handles GET requests to retrieve a single user
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
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
