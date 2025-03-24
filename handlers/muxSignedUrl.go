package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
)

func HandleMuxSignedUploadUrl(c *gin.Context) {
	muxSignedUrlResponse, err := internal.Mux()

	if err != nil {
		log.Println("ERROR HANDLING SIGNED URL: ", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, {url: muxSignedUrlResponse.Data.Url})
}
