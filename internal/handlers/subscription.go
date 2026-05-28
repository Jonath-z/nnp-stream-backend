package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
	"github.com/nnp-stream-backend/models"
)

// HandleSubscribeToPlan godoc
// @Summary Subscribe to a plan via Shwary mobile money (DRC, USD)
// @Description Initiates a Shwary mobile-money payment for the given plan and creates a pending subscription. The subscription is activated when the Shwary webhook reports a completed transaction.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param payload body models.SubscribeRequest true "Subscribe payload"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func HandleSubscribeToPlan(c *gin.Context) {
	var req models.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phone := strings.TrimSpace(req.PhoneNumber)
	if !strings.HasPrefix(phone, "+") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone_number must be in E.164 format (e.g. +243...)"})
		return
	}
	if !strings.HasPrefix(phone, "+243") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only DRC numbers (+243) are supported"})
		return
	}

	plan, err := internal.GetPlan(req.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan not found"})
		return
	}

	subscription, err := internal.CreateSubscription(models.Subscription{
		UserID: req.UserID,
		PlanID: plan.ID,
		Status: models.SubscriptionStatusPending,
	})
	if err != nil {
		log.Printf("subscribe: create subscription failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	shwaryTx, err := internal.InitiateShwaryPayment(req.Amount, phone, subscription.ID)
	if err != nil {
		log.Printf("subscribe: shwary payment failed: %v", err)
		_ = internal.UpdateSubscriptionStatus(subscription.ID, models.SubscriptionStatusFailed, "")
		c.JSON(http.StatusBadGateway, gin.H{"error": "payment provider error"})
		return
	}

	if _, err := internal.CreatePaymentTransaction(models.PaymentTransaction{
		UserID:         req.UserID,
		PlanID:         plan.ID,
		SubscriptionID: subscription.ID,
		Provider:       "shwary",
		ProviderTxID:   shwaryTx.ID,
		Amount:         req.Amount,
		Currency:       "USD",
		PhoneNumber:    phone,
		Status:         shwaryTx.Status,
	}); err != nil {
		log.Printf("subscribe: persist transaction failed: %v", err)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"subscription": subscription,
		"transaction":  shwaryTx,
	})
}
