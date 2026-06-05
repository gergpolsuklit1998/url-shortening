package handlers

import (
	"context"
	"net/http"
	_ "strconv"

	"github.com/segmentio/ksuid"

	"github.com/gergpolsuklit1998/url-shortening/models"
	"github.com/gergpolsuklit1998/url-shortening/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

// ShortUrlHandler contains all HTTP handlers for short url operations.
// It depends on the repository for database access.
type ShortUrlHandler struct {
	repo *repository.ShortUrlRepository
}

// NewShortUrlHandler creates a new handler with the given repository.
func NewShortUrlHandler(repo *repository.ShortUrlRepository) *ShortUrlHandler {
	return &ShortUrlHandler{
		repo: repo,
	}
}

func (h *ShortUrlHandler) increaseAccessCount(ctx context.Context, shortUrl *models.ShortUrl) (*models.ShortUrl, error) {
	shortUrl.AccessCount++
	return h.repo.UpdateAccessCount(ctx, shortUrl.ID, shortUrl.AccessCount)
}

// CreateShortUrl handles POST /shortens requests.
// It validates the request body and creates a new short url.
func (h *ShortUrlHandler) CreateShortUrl(c *gin.Context) {
	var req models.CreateShortUrlRequest

	// Bind and validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"error":      "Invalid request body",
		})
		return
	}

	// Generate shortcode
	shortCode := ksuid.New().String()[:6]

	// Create short url model from request
	shortUrl := &models.ShortUrl{
		Url:         req.Url,
		ShortCode:   shortCode,
		AccessCount: 0,
	}

	// Save to database
	createdShorUrl, err := h.repo.Create(c.Request.Context(), shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Internal Server Error",
		})
		return
	}

	createdShorUrl.ShortCode = "http://localhost:8000/api/v1/shortens/" + shortCode
	// Return created short url
	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    "Short url created successfully",
		"data":       createdShorUrl,
	})
}

// RedirectShortUrl handles GET /shortens/:shortCode requests.
// It redirects to the original URL.
func (h *ShortUrlHandler) RedirectShortUrl(c *gin.Context) {
	// Parse the ShortCode parameter from the URL
	shortCodeParam := c.Param("shortCode")

	// Fetch short url from database
	shorUrl, err := h.repo.FindByShortCode(c.Request.Context(), shortCodeParam)
	if err != nil {
		if err == repository.ErrShortUrlNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"error":      repository.ErrShortUrlNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Internal Server Error",
		})
		return
	}

	// Increase accessCounter
	_, err = h.increaseAccessCount(c.Request.Context(), shorUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Internal Server Error",
		})
		return
	}

	// Redirect to original URL
	c.Redirect(http.StatusMovedPermanently, shorUrl.Url)
}

func (h *ShortUrlHandler) UpdateShortUrl(c *gin.Context) {
	// Parse the ShortCode parameter from the URL
	shortCodeParam := c.Param("shortCode")

	// Bind the update request
	var req models.UpdateShortUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"error":      "Invalid request body",
		})
		return
	}

	// Build update document with only provided fields
	update := bson.M{}
	if req.Url != nil {
		update["url"] = *req.Url
	}

	// Check if there is anything to update
	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"error":      "No fields to update",
		})
		return
	}

	// Perform the update
	updatedShortUrl, err := h.repo.Update(c.Request.Context(), shortCodeParam, update)
	if err != nil {
		if err == repository.ErrShortUrlNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"error":      repository.ErrShortUrlNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Internal Server Error",
		})
		return
	}

	updatedShortUrl.ShortCode = "http://localhost:8000/api/v1/shortens/" + updatedShortUrl.ShortCode

	c.JSON(http.StatusOK, gin.H{
		"message": "Short url updated successfully",
		"data":    updatedShortUrl,
	})
}

func (h *ShortUrlHandler) DeleteShortUrl(c *gin.Context) {
	shortCodeParam := c.Param("shortCode")

	err := h.repo.Delete(c.Request.Context(), shortCodeParam)
	if err != nil {
		if err == repository.ErrShortUrlNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"error":      repository.ErrShortUrlNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Short url deleted successfully",
	})
}
