package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gergpolsuklit1998/url-shortening/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Common errors returned by the repository
var (
	ErrShortUrlNotFound = errors.New("short url not found")
)

// ShortUrlRepository handles all short url-related database operations.
// It abstracts MongoDB operations behind a clean interface.
type ShortUrlRepository struct {
	collection *mongo.Collection
}

// NewShortUrlRepository creates a new repository instance.
// Pass the MongoDB collection to use for short url documents.
func NewShortUrlRepository(collection *mongo.Collection) *ShortUrlRepository {
	return &ShortUrlRepository{
		collection: collection,
	}
}

// Create inserts a new short url document into the database.
// It sets timestamps and returns the created short url with its ID.
func (r *ShortUrlRepository) Create(ctx context.Context, shortUrl *models.ShortUrl) (*models.ShortUrl, error) {
	// Set timestamps
	now := time.Now()
	shortUrl.CreatedAt = now
	shortUrl.UpdatedAt = now

	// Insert the document
	result, err := r.collection.InsertOne(ctx, shortUrl)
	if err != nil {
		return nil, err
	}

	// Set the generated ID on the short url struct
	shortUrl.ID = result.InsertedID.(primitive.ObjectID)

	return shortUrl, nil
}

func (r *ShortUrlRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.ShortUrl, error) {
	var shortUrl models.ShortUrl

	// Create filter to match the short_code field
	filter := bson.M{"short_code": shortCode}

	// Execute the query
	err := r.collection.FindOne(ctx, filter).Decode(&shortUrl)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrShortUrlNotFound
		}
		return nil, err
	}

	return &shortUrl, nil
}

func (r *ShortUrlRepository) UpdateAccessCount(ctx context.Context, id primitive.ObjectID, accessCount int) (*models.ShortUrl, error) {
	// Create filter to match the _id field
	filter := bson.M{"_id": id}

	// Create update document
	update := bson.M{"$set": bson.M{"access_count": accessCount}}

	// Execute the update
	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var shortUrl models.ShortUrl
	err := result.Decode(&shortUrl)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrShortUrlNotFound
		}
		return nil, err
	}

	return &shortUrl, nil
}
