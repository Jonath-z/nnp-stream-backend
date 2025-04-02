package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
)

func HandleMuxSignedUploadUrl(c *gin.Context) {
	muxSignedUrlResponse, err := internal.GenerateMuxUploadURL()

	if err != nil {
		log.Print("ERROR HANDLING SIGNED URL: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, muxSignedUrlResponse.Data.Url)
}
