package webhook

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
	"github.com/nnp-stream-backend/models"
)

// HandleShwaryWebhook godoc
// @Summary Shwary callback receiver
// @Description Receives Shwary transaction status callbacks (pending, completed, failed, cancelled) and updates the matching subscription/transaction. Endpoint is idempotent.
// @Tags webhooks
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /shwary-web-hook [post]
func HandleShwaryWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("shwary webhook: read body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var payload models.ShwaryTransaction
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("shwary webhook: unmarshal: %v body=%s", err, string(body))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if payload.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing transaction id"})
		return
	}

	storedTx, err := internal.GetPaymentTransactionByProviderID("shwary", payload.ID)
	if err != nil {
		log.Printf("shwary webhook: transaction %s not found: %v", payload.ID, err)
		// Ack so Shwary doesn't retry, but log it. Reconciliation will catch it.
		c.JSON(http.StatusOK, gin.H{"message": "unknown transaction acknowledged"})
		return
	}

	patch := map[string]interface{}{
		"status": payload.Status,
	}
	if payload.FailureReason != "" {
		patch["failure_reason"] = payload.FailureReason
	}
	if payload.TxHash != "" {
		patch["tx_hash"] = payload.TxHash
	}
	if payload.CompletedAt != "" {
		patch["completed_at"] = payload.CompletedAt
	}
	if err := internal.UpdatePaymentTransaction(storedTx.ID, patch); err != nil {
		log.Printf("shwary webhook: update transaction: %v", err)
	}

	if storedTx.SubscriptionID != "" {
		switch payload.Status {
		case "completed":
			plan, planErr := internal.GetPlan(storedTx.PlanID)
			expiresAt := ""
			if planErr == nil {
				expiresAt = internal.ComputeSubscriptionExpiry(plan.BillingCycle, time.Now().UTC())
			}
			if err := internal.UpdateSubscriptionStatus(storedTx.SubscriptionID, models.SubscriptionStatusActive, expiresAt); err != nil {
				log.Printf("shwary webhook: activate subscription: %v", err)
			}
		case "failed":
			if err := internal.UpdateSubscriptionStatus(storedTx.SubscriptionID, models.SubscriptionStatusFailed, ""); err != nil {
				log.Printf("shwary webhook: fail subscription: %v", err)
			}
		case "cancelled":
			if err := internal.UpdateSubscriptionStatus(storedTx.SubscriptionID, models.SubscriptionStatusCanceled, ""); err != nil {
				log.Printf("shwary webhook: cancel subscription: %v", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
