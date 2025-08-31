package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GET /:shortCode
func (h *URLHandler) RedirectToURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	originalURL, err := h.urlService.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// 301 status code with original URL in Location header
	c.Redirect(http.StatusMovedPermanently, originalURL)
}
