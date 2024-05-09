package repository

import (
	"context"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Interface for querying posts from db
type PostRepository interface {
	CreatePost(ctx context.Context, post model.Post) error
	GetPosts(ctx context.Context) ([]model.Post, error)
}

type postRepository struct {
	db *mongo.Collection
}

// Create a new post repository
func NewPostRepository(db *mongo.Database) PostRepository {
	return &postRepository{
		db: db.Collection("posts"),
	}
}
// Inserts a new post into the database
func (r *postRepository) CreatePost(ctx context.Context, post model.Post) error {
	_, err := r.db.InsertOne(ctx, post)
    return err
}

// FindAll returns all posts from the database
func (r *postRepository) GetPosts(ctx context.Context) ([]model.Post, error) {
    var posts []model.Post
    cur, err := r.db.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)
    for cur.Next(ctx) {
        var post model.Post
        if err := cur.Decode(&post); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    if err := cur.Err(); err != nil {
        return nil, err
    }
    return posts, nil
}

