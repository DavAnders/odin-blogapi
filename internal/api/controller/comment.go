package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"github.com/DavAnders/odin-blogapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController struct {
	repo repository.CommentRepository
}

func NewCommentController(repo repository.CommentRepository) *CommentController {
	return &CommentController{
		repo: repo,
	}
}

// Handles POST requests to create a new comment
func (c *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid comment data", http.StatusBadRequest)
		return
	}
	comment.ID = primitive.NewObjectID()

	if err := c.repo.CreateComment(context.Background(), comment); err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// Handles GET requests to retrieve all comments for a post
func (c *CommentController) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("postId")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	comments, err := c.repo.GetCommentsByPost(context.Background(), postID)
	if err != nil {
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
