package handlers

import (
	"net/http"
	_ "strconv"

	"github.com/gergpolsuklit1998/url-shortening/models"
	"github.com/gergpolsuklit1998/url-shortening/repository"
	"github.com/gergpolsuklit1998/url-shortening/services"

	"github.com/gin-gonic/gin"
)

// ShortUrlHandler contains all HTTP handlers for short url operations.
// It depends on the repository for database access.
type ShortUrlHandler struct {
	service *services.ShortUrlService
}

// NewShortUrlHandler creates a new handler with the given repository.
func NewShortUrlHandler(service *services.ShortUrlService) *ShortUrlHandler {
	return &ShortUrlHandler{
		service: service,
	}
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

	createdShorUrl, err := h.service.CreateShortUrl(c.Request.Context(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Internal Server Error",
		})
		return
	}

	createdShorUrl.ShortCode = "http://localhost:8000/api/v1/shortens/" + createdShorUrl.ShortCode
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
	shorUrl, err := h.service.GetShortUrl(c.Request.Context(), shortCodeParam)
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
	_, err = h.service.IncreaseAccessCount(c.Request.Context(), shorUrl)
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

func (h *ShortUrlHandler) GetShortUrlStats(c *gin.Context) {
	shortCodeParam := c.Param("shortCode")

	shortUrl, err := h.service.GetShortUrl(c.Request.Context(), shortCodeParam)
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
		"message":    "Short url stats retrieved successfully",
		"data":       shortUrl,
	})
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

	updatedShortUrl, err := h.service.UpdateShortUrl(c.Request.Context(), shortCodeParam, &req)
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
		"statusCode": http.StatusOK,
		"message":    "Short url updated successfully",
		"data":       updatedShortUrl,
	})
}

func (h *ShortUrlHandler) DeleteShortUrl(c *gin.Context) {
	shortCodeParam := c.Param("shortCode")

	err := h.service.DeleteShortUrl(c.Request.Context(), shortCodeParam)
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
