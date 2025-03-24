package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/models"
)

type Handler struct {
}

func HandleMuxWebhook(c *gin.Context) {
	log.Println(c.Request.URL)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error(), "ERRRRRRO RRRRRRR")
	}

	var muxWebhookPayload models.MuxWebhookPayload
	marchalErr := json.Unmarshal(body, muxWebhookPayload)
	if marchalErr != nil {
		log.Print("Unable to unmarshal webhook payload")
	}
	log.Println("Data", string(body), muxWebhookPayload.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook received successfully",
	})
}
