package repository

import (
	"context"
	"fmt"
	"log"
	"time"

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
	if user.Email == "" {
        return fmt.Errorf("email is required")
    }
    if user.Username == "" {
        return fmt.Errorf("username is required")
    }
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = "" // Clear the plain password
	user.CreatedAt = time.Now()

	if user.ID.IsZero() {
        user.ID = primitive.NewObjectID()
    }

	result, err := r.db.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}
	log.Printf("Inserted user with ID: %v", result.InsertedID)
	
	return nil
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