package repository

import (
	"context"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment model.Comment) error
	GetCommentsByPost(ctx context.Context, postID string) ([]model.Comment, error)
}

type commentRepository struct {
	db *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) CommentRepository {
	return &commentRepository{
		db: db.Collection("comments"),
	}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment model.Comment) error {
	_, err := r.db.InsertOne(ctx, comment)
	return err
}

func (r *commentRepository) GetCommentsByPost(ctx context.Context, postID string) ([]model.Comment, error) {
	var comments []model.Comment
	filter := bson.M{"postId": postID}
	cur, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var comment model.Comment
		if err := cur.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
