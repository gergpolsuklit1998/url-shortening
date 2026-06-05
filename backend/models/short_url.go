package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShortUrl represents a short URL document in MongoDB.
// The bson tags define how fields are stored in MongoDB.
// The json tags define how fields appear in API responses.
type ShortUrl struct {
	// ID is the MongoDB ObjectID, automatically generated if not provided
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Url string `bson:"url" json:"url" binding:"required"`

	ShortCode string `bson:"short_code" json:"short_code"`

	AccessCount int `bson:"access_count" json:"access_count"`

	// CreatedAt is set automatically when the document is created
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	// UpdatedAt is updated whenever the document changes
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// CreateShortUrlRequest contains the fields required to create a new short URL.
// Using a separate struct prevents clients from setting internal fields.
type CreateShortUrlRequest struct {
	Url         string `json:"url" binding:"required"`
	ShortCode   string `json:"short_code"`
	AccessCount int    `json:"access_count"`
}

// UpdateShortUrlRequest contains fields that can be updated.
// All fields are optional to support partial updates.
type UpdateShortUrlRequest struct {
	Url       *string `json:"url" binding:"omitempty"`
	ShortCode *string `json:"short_code" binding:"omitempty"`
}
