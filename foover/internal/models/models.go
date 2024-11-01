package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session represents a user session
type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SessionID string             `bson:"session_id"`
	CreatedAt time.Time          `bson:"created_at"`
}

// Vote represents a vote on a product by a session
type Vote struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SessionID string             `bson:"session_id"`
	ProductID string             `bson:"product_id"`
	Score     int                `bson:"score"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// ProductScore represents the aggregated score of a product
type ProductScore struct {
	ProductID string  `bson:"_id"`
	AvgScore  float64 `bson:"avg_score"`
	VoteCount int     `bson:"vote_count"`
}

// Product represents a product with an ID
type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ProductID string             `bson:"product_id" json:"product_id"`
}
