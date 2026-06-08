package services

import (
	"context"
	"errors"

	"github.com/gergpolsuklit1998/url-shortening/models"
	"github.com/gergpolsuklit1998/url-shortening/repository"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
)

type ShortUrlService struct {
	repo *repository.ShortUrlRepository
}

func NewShortUrlService(repo *repository.ShortUrlRepository) *ShortUrlService {
	return &ShortUrlService{
		repo: repo,
	}
}

func (s *ShortUrlService) IncreaseAccessCount(ctx context.Context, shortUrl *models.ShortUrl) (*models.ShortUrl, error) {
	shortUrl.AccessCount++
	return s.repo.UpdateAccessCount(ctx, shortUrl.ID, shortUrl.AccessCount)
}

func (s *ShortUrlService) CreateShortUrl(ctx context.Context, req *models.CreateShortUrlRequest) (*models.ShortUrl, error) {
	// Generate shortcode
	shortCode := ksuid.New().String()[:6]

	// Create short url model from request
	shortUrl := &models.ShortUrl{
		Url:         req.Url,
		ShortCode:   shortCode,
		AccessCount: 0,
	}

	// Save to database
	return s.repo.Create(ctx, shortUrl)
}

func (s *ShortUrlService) GetShortUrl(ctx context.Context, shortCode string) (*models.ShortUrl, error) {
	return s.repo.FindByShortCode(ctx, shortCode)
}

func (s *ShortUrlService) UpdateShortUrl(ctx context.Context, shortCode string, req *models.UpdateShortUrlRequest) (*models.ShortUrl, error) {
	// Build update document with only provided fields
	update := bson.M{}
	if req.Url != nil {
		update["url"] = *req.Url
	}

	// Check if there is anything to update
	if len(update) == 0 {
		return nil, errors.New("no fields to update")
	}

	// Perform the update
	return s.repo.Update(ctx, shortCode, update)
}

func (s *ShortUrlService) DeleteShortUrl(ctx context.Context, shortCode string) error {
	return s.repo.Delete(ctx, shortCode)
}
