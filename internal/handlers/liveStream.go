package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
)

func HandleMuxLiveStream(c *gin.Context) {
	liveStream, err := internal.CreateLiveStream()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, liveStream)
}

func HandleMuxDeleteLiveStream(c *gin.Context) {
	liveStreamID := c.Param("liveStreamID")
	err := internal.DeleteLiveStream(liveStreamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Live stream deleted successfully"})
}