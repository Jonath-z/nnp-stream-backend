package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Health check
// @Description Returns a welcome message to verify the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to nnp stream backend service",
	})
}
