package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
)

// HandleMuxLiveStream godoc
// @Summary Create a live stream
// @Description Creates a new low-latency Mux live stream with public playback policy
// @Tags live-streams
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /mux-live-stream [post]
func HandleMuxLiveStream(c *gin.Context) {
	liveStream, err := internal.CreateLiveStream()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, liveStream)
}

// HandleMuxDeleteLiveStream godoc
// @Summary Delete a live stream
// @Description Deletes a Mux live stream by its ID
// @Tags live-streams
// @Produce json
// @Param liveStreamID path string true "Live Stream ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /mux-live-stream/{liveStreamID} [delete]
func HandleMuxDeleteLiveStream(c *gin.Context) {
	liveStreamID := c.Param("liveStreamID")
	err := internal.DeleteLiveStream(liveStreamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Live stream deleted successfully"})
}