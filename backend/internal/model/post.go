package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title string `bson:"title" json:"title" binding:"required"`
	Content string `bson:"content" json:"content" binding:"required"`
	Published bool `bson:"published" json:"published"`
	PublishedAt time.Time `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	AuthorID primitive.ObjectID `bson:"authorId" json:"authorId"`
}