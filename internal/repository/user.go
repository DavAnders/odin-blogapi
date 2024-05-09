package repository

import (
	"context"

	"github.com/DavAnders/odin-blogapi/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	ValidateCredentials(ctx context.Context, username, password string) (*model.User, error)
}

type userRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db.Collection("users"),
	}
}

// Inserts a new user into the database
func (r *userRepository) CreateUser(ctx context.Context, user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = "" // Clear the plain password
	user.ID = primitive.NewObjectID() // Ensure an ObjectID is generated

	_, err = r.db.InsertOne(ctx, user)
	return err
}

// Returns a single user from the database
func (r *userRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Returns all users from the database
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

// Validates a user's credentials
func (r *userRepository) ValidateCredentials(ctx context.Context, username, password string) (*model.User, error) {
    var user model.User
    err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        return nil, err
    }

    // Compare the stored hashed password with the provided password
    if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
        return nil, err  // Password does not match
    }

    return &user, nil
}