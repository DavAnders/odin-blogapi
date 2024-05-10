package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DavAnders/odin-blogapi/backend/internal/model"
	"github.com/DavAnders/odin-blogapi/backend/internal/repository"
	"github.com/go-chi/chi/v5"
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
    postID := chi.URLParam(r, "id")
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

// Handles PUT requests to update a comment
func (c *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "id")
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

// Handles DELETE requests to delete a comment
func (c *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "id")
    if commentID == "" {
        http.Error(w, "Comment ID is required", http.StatusBadRequest)
        return
    }

    // Extract userID from context, make sure it's available
    userIDValue := r.Context().Value("userID")  // Retrieve the userID from context
    userID, ok := userIDValue.(string)
    if !ok || userID == "" {
        http.Error(w, "Unauthorized or bad request", http.StatusUnauthorized)
        return
    }

    // Convert userID to pointer for repository call
    userIDPtr := &userID

    // Delete the comment directly with user authorization check in the repo layer
    if err := c.repo.DeleteComment(context.Background(), commentID, userIDPtr); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // No Content is typical for a successful delete operation
}


// Deletes a comment by ID, only if the user is an admin
func (c *CommentController) AdminDeleteComment(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "id")
    if commentID == "" {
        http.Error(w, "Comment ID is required", http.StatusBadRequest)
        return
    }

    // Pass nil as userID to indicate an admin deletion
    if err := c.repo.DeleteComment(context.Background(), commentID, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

