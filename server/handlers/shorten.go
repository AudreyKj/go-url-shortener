package handlers

import (
	"go-url-shortner/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /api/urls
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req services.URLRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validate and normalize the input URL
	validator := services.NewURLValidator()
	normalizedURL, err := validator.NormalizeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}
	// Use normalized URL for further processing
	req.URL = normalizedURL

	response, err := h.urlService.CreateShortURL(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
