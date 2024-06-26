package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/DavAnders/odin-blogapi/backend/internal/api/controller"
	"github.com/DavAnders/odin-blogapi/backend/internal/api/middleware"
	"github.com/DavAnders/odin-blogapi/backend/internal/repository"
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
	postRepo := repository.NewPostRepository(client.Database("blogprod"))
	userRepo := repository.NewUserRepository(client.Database("blogprod"))
	commentRepo := repository.NewCommentRepository(client.Database("blogprod"))

	// Initialize controllers
	postController := controller.NewPostController(postRepo)
	userController := controller.NewUserController(userRepo)
	commentController := controller.NewCommentController(commentRepo)

	r := chi.NewRouter()

	// Serve files
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL Path: %s", r.URL.Path)
		path := filepath.Join("public", r.URL.Path)
	
		// Check if the file exists and is not a directory
		if stat, err := os.Stat(path); os.IsNotExist(err) || stat.IsDir() {
			log.Println("File does not exist or is a directory, serving index.html")
			http.ServeFile(w, r, "public/index.html")
		} else {
			log.Printf("Serving static file: %s", path)
			http.ServeFile(w, r, path)
		}
	})

	// Apply CORS middleware
	r.Use(middleware.EnableCORS)

	// Public routes
	r.Post("/login", userController.Login)
	r.Post("/register", userController.Register)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // Apply auth middleware to all '/api' routes

		r.Get("/posts", postController.GetPosts)
		r.Post("/posts", postController.CreatePost)
		r.Get("/posts/user/{userID}", postController.GetPostsByUser)
		r.Post("/posts", postController.CreatePost)
		r.Get("/posts/{id}", postController.GetPostByID)
		r.Put("/posts/{id}", postController.UpdatePost)
		r.Delete("/posts/{id}", postController.DeletePost)

		r.Get("/profile", userController.GetUserProfile)
        r.Put("/profile", userController.UpdateUserProfile)

		r.Get("/users", userController.GetUsers)
		r.Post("/users", userController.CreateUser)
		r.Get("/users/{id}", userController.GetUser)

		r.Post("/comments", commentController.CreateComment)
		r.Get("/comments/{id}", commentController.GetCommentsByPost)
		r.Put("/comments/{id}", commentController.UpdateComment)
		r.Delete("/comments/{id}", commentController.DeleteComment)

		// Admin-specific routes under '/api/admin'
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminMiddleware(*repository.NewAdminRepository(client.Database("blog")))) // Apply admin-specific middleware
			r.Delete("/posts/{id}", postController.AdminDeletePost)
			r.Delete("/comments/{id}", commentController.AdminDeleteComment)
		})
	})
	
	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	

	// Start server
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
