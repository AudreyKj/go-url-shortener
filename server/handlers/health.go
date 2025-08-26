package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /health requests
func (h *URLHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"storage": "redis",
	})
}
