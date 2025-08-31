package handlers

import (
	"net/http"
	"go-url-shortner/services"
	"github.com/gin-gonic/gin"
)

// POST /api/urls 
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req services.URLRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	response, err := h.urlService.CreateShortURL(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
