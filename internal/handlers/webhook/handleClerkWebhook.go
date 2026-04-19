package webhook

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
	"github.com/nnp-stream-backend/models/event"
)

// HandleCleckWebhook godoc
// @Summary Clerk webhook receiver
// @Description Receives Clerk webhook events (e.g. user.created) and syncs users
// @Tags webhooks
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /clerk-web-hook [post]
func HandleCleckWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error(), "ERROR")
	}
	var clerkWebhookPayload event.ClerkUserEvent
	unmarchallErr := json.Unmarshal(body, &clerkWebhookPayload)
	if unmarchallErr != nil {
		fmt.Println("Error while unmarchalling", unmarchallErr)
	}

	if clerkWebhookPayload.Type == "user.created" {
		internal.CreateUser(clerkWebhookPayload)
	}

	c.JSON(200, gin.H{
		"message": "Webhook received successfully",
	})

}
