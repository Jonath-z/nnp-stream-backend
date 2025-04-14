package webhook

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
	"github.com/nnp-stream-backend/models/event"
)

func HandleCleckWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error(), "ERROR")
	}
	var clerkWebhookPayload event.ClerkUserEvent
	json.Unmarshal(body, &clerkWebhookPayload)

	if clerkWebhookPayload.Type == "user.created" {
		internal.CreateUser(clerkWebhookPayload)
	}

	c.JSON(200, gin.H{
		"message": "Webhook received successfully",
	})

}
