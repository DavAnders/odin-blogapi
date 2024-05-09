package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"github.com/DavAnders/odin-blogapi/internal/repository"

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