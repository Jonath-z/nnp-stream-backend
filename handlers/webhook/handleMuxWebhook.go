package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
	"github.com/nnp-stream-backend/models"
)

func HandleMuxWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error(), "ERROR")
	}

	var muxWebhookPayload models.VideoAssetReadyEvent
	marchalErr := json.Unmarshal(body, &muxWebhookPayload)
	if marchalErr != nil {
		log.Print("Unable to unmarshal webhook payload", marchalErr.Error())
	}

	log.Println("Data", muxWebhookPayload.ID)

	if muxWebhookPayload.Type == "video.asset.ready" {
		internal.UploadDraftAsset(models.UploadDraftPayload{
			AssetID:    muxWebhookPayload.Data.UploadID,
			PlaybackId: muxWebhookPayload.Data.PlaybackIDs[0].ID,
			Duration:   muxWebhookPayload.Data.Duration,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook received successfully",
	})
}
