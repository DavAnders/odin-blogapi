package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Interface for querying posts from db
type PostRepository interface {
	CreatePost(ctx context.Context, post model.Post) error
	GetPosts(ctx context.Context) ([]model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	UpdatePost(ctx context.Context, id string, post model.Post) error
	DeletePost(ctx context.Context, id string) error
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

// Find a post by its ID
func (r *postRepository) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
    var post model.Post
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err  // If the ID is not a valid ObjectId
    }
    if err := r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&post); err != nil {
        return nil, err
    }
    return &post, nil
}

// Updates a post in the database
func (r *postRepository) UpdatePost(ctx context.Context, id string, post model.Post) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err  // If the ID is not a valid ObjectId
    }
    update := bson.M{
        "$set": bson.M{
            "title": post.Title,
            "content": post.Content,
            "published": post.Published,
            "updatedAt": time.Now(),
        },
    }
    filter := bson.M{"_id": objID}
    result, err := r.db.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
		return fmt.Errorf("no post found with given ID")
    }
    return nil
}

// Deletes a post from the database
func (r *postRepository) DeletePost(ctx context.Context, id string) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err  // If the ID is not a valid ObjectId
    }
    filter := bson.M{"_id": objID}
    result, err := r.db.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return fmt.Errorf("no post found with given ID")
    }
    return nil
}
