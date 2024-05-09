package repository

import (
	"context"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context) ([]model.User, error)
}

type userRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db.Collection("users"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.db.InsertOne(ctx, user)
	return err
}

func (r *userRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	cur, err := r.db.Find(ctx, bson.M{}, options.Find())  // Can limit / sort if needed
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user model.User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return users, nil
}