package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username" binding:"required"`
	Password string `bson:"-" json:"password"`
	HashedPassword string `bson:"hashedPassword" json:"-"`
	Email string `bson:"email" json:"email" binding:"required"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}