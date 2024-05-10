package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"github.com/DavAnders/odin-blogapi/internal/repository"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostController struct {
	repo repository.PostRepository
}

func NewPostController(repo repository.PostRepository) *PostController {
	return &PostController{
		repo: repo,
	}
}

// Handles POST requests
func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.ID = primitive.NewObjectID()
	post.PublishedAt = time.Now()

	if err := c.repo.CreatePost(context.Background(), post); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// Handles GET requests to retrieve all posts
func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
    posts, err := c.repo.GetPosts(context.Background())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

// Handles GET requests to retrieve a post by ID
func (c *PostController) GetPostByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    post, err := c.repo.GetPostByID(context.Background(), postID)
    if err != nil {
        http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

// Handles PUT requests to update a post
func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    var post model.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := c.repo.UpdatePost(context.Background(), postID, post); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Explicitly signify a successful update
    json.NewEncoder(w).Encode(post)
}

// Handles DELETE requests to delete a post
func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]
    userID := r.Context().Value("userID").(string) // Get user ID from context, cast to string

    post, err := c.repo.GetPostByID(context.Background(), postID)
    if err != nil {
        http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
        return
    }
    if post == nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    if post.AuthorID.Hex() != userID {
        http.Error(w, "Unauthorized to delete this post", http.StatusUnauthorized)
        return
    }

    // Proceed to delete the post
    if err := c.repo.DeletePost(context.Background(), postID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
