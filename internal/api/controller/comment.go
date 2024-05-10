package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"github.com/DavAnders/odin-blogapi/internal/repository"
	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
    postID := vars["id"]
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

func (c *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    commentID := vars["id"]
    if commentID == "" {
        http.Error(w, "Comment ID is required", http.StatusBadRequest)
        return
    }

    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized or bad request", http.StatusUnauthorized)
        return
    }

    var comment model.Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Update the comment directly with user authorization check in the repo layer
    if err := c.repo.UpdateComment(context.Background(), commentID, userID, comment); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Explicitly signify a successful update
    json.NewEncoder(w).Encode(comment)
}
