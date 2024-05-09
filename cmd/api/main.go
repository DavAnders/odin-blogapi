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
	postRepo := repository.NewPostRepository(client.Database("yourDatabaseName"))
	userRepo := repository.NewUserRepository(client.Database("yourDatabaseName"))
	commentRepo := repository.NewCommentRepository(client.Database("yourDatabaseName"))

	// Initialize controllers
	postController := controller.NewPostController(postRepo)
	userController := controller.NewUserController(userRepo)
	commentController := controller.NewCommentController(commentRepo)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/posts", postController.GetPosts).Methods("GET")
	r.HandleFunc("/posts", postController.CreatePost).Methods("POST")
	r.HandleFunc("/users", userController.GetUsers).Methods("GET")
	r.HandleFunc("/users", userController.CreateUser).Methods("POST")
	r.HandleFunc("/comments", commentController.GetCommentsByPost).Methods("GET")
	r.HandleFunc("/comments", commentController.CreateComment).Methods("POST")

	// Apply middleware
	r.Use(middleware.AuthMiddleware)

	// Start server
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
