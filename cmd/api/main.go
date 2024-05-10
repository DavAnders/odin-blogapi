package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/DavAnders/odin-blogapi/internal/api/controller"
	"github.com/DavAnders/odin-blogapi/internal/api/middleware"
	"github.com/DavAnders/odin-blogapi/internal/repository"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set in .env file")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Connected to MongoDB")

	// Initialize repositories
	postRepo := repository.NewPostRepository(client.Database("blog"))
	userRepo := repository.NewUserRepository(client.Database("blog"))
	commentRepo := repository.NewCommentRepository(client.Database("blog"))

	// Initialize controllers
	postController := controller.NewPostController(postRepo)
	userController := controller.NewUserController(userRepo)
	commentController := controller.NewCommentController(commentRepo)

	// Main router
	r := mux.NewRouter()

	// Subrouter for authenticated routes
	authRoutes := r.PathPrefix("/api").Subrouter()
	authRoutes.Use(middleware.AuthMiddleware)  // Apply middleware to all routes in this subrouter

	// Admin routes for wrapping the admin middleware
	adminRepo := repository.NewAdminRepository(client.Database("blog"))

	adminRoutes := r.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.AuthMiddleware) // Ensures the user is authenticated
	adminRoutes.Use(middleware.AdminMiddleware(*adminRepo)) // Ensures the user is an admin


	// Public routes
	r.HandleFunc("/login", userController.Login).Methods("POST")
	r.HandleFunc("/register", userController.Register).Methods("POST")  // For user registration with JWT token

	// Authenticated routes
	authRoutes.HandleFunc("/posts", postController.GetPosts).Methods("GET")
	authRoutes.HandleFunc("/posts", postController.CreatePost).Methods("POST") 
	authRoutes.HandleFunc("/posts/{id}", postController.GetPostByID).Methods("GET")
	authRoutes.HandleFunc("/posts/{id}", postController.UpdatePost).Methods("PUT")
	authRoutes.HandleFunc("/posts/{id}", postController.DeletePost).Methods("DELETE")
	authRoutes.HandleFunc("/users", userController.GetUsers).Methods("GET")
	authRoutes.HandleFunc("/users", userController.CreateUser).Methods("POST") // For creating a user without JWT (admin usage)
	authRoutes.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")
	authRoutes.HandleFunc("/comments/{id}", commentController.GetCommentsByPost).Methods("GET")
	authRoutes.HandleFunc("/comments/{id}", commentController.UpdateComment).Methods("PUT")
	authRoutes.HandleFunc("/comments/{id}", commentController.DeleteComment).Methods("DELETE")
	authRoutes.HandleFunc("/comments", commentController.CreateComment).Methods("POST")

	// Start server
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
