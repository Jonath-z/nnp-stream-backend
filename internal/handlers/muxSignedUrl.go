package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
)

// HandleMuxSignedUploadUrl godoc
// @Summary Get a signed upload URL
// @Description Generates a new Mux direct upload URL for uploading video assets
// @Tags uploads
// @Produce json
// @Success 200 {string} string "Signed upload URL"
// @Failure 400 {object} map[string]string
// @Router /mux-signed-url [get]
func HandleMuxSignedUploadUrl(c *gin.Context) {
	muxSignedUrlResponse, err := internal.GenerateMuxUploadURL()

	if err != nil {
		log.Print("ERROR HANDLING SIGNED URL: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, muxSignedUrlResponse.Data.Url)
}
