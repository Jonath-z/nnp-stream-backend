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

	user, userError := internal.GetUserByClerkUserId(req.UserID)
	if userError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
	}

	subscription, err := internal.CreateSubscription(models.Subscription{
		UserID: user.Id,
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
		UserID:         user.Id,
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

// GetSubscriptionStatus godoc
// @Summary Check whether a user has access to a video
// @Description Looks up the plans that gate the given video and returns whether the user has an active, unexpired subscription to any of them. If no plan gates the video, access is granted. If only a pending subscription exists, has_access is false with reason "pending".
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param payload body models.CheckVideoAccessRequest true "Access check payload"
// @Success 200 {object} models.CheckVideoAccessResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/access [post]
func GetSubscriptionStatus(c *gin.Context) {
	var req models.CheckVideoAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := internal.GetUserByClerkUserId(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	planIDs, err := internal.GetPlanIDsForVideo(req.VideoID)
	if err != nil {
		log.Printf("access check: load video plans: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve video access"})
		return
	}

	if len(planIDs) == 0 {
		c.JSON(http.StatusOK, models.CheckVideoAccessResponse{
			HasAccess: true,
			Reason:    models.VideoAccessReasonFreeVideo,
		})
		return
	}

	sub, err := internal.FindUserSubscriptionForPlans(user.Id, planIDs)
	if err != nil {
		log.Printf("access check: load subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve subscription"})
		return
	}

	if sub == nil {
		c.JSON(http.StatusOK, models.CheckVideoAccessResponse{
			HasAccess: false,
			Reason:    models.VideoAccessReasonNoSubscription,
		})
		return
	}

	resp := models.CheckVideoAccessResponse{Subscription: sub}
	switch sub.Status {
	case models.SubscriptionStatusActive:
		resp.HasAccess = true
		resp.Reason = models.VideoAccessReasonActive
	case models.SubscriptionStatusPending:
		resp.HasAccess = false
		resp.Reason = models.VideoAccessReasonPending
	default:
		resp.HasAccess = false
		resp.Reason = models.VideoAccessReasonExpired
	}
	c.JSON(http.StatusOK, resp)
}

// PollSubscriptionStatus godoc
// @Summary Poll a subscription's status
// @Description Returns the current state of a subscription. Used by the frontend to poll after initiating a Shwary payment until the subscription transitions out of "pending".
// @Tags subscriptions
// @Produce json
// @Param subscriptionID path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{subscriptionID} [get]
func PollSubscriptionStatus(c *gin.Context) {
	subscriptionID := strings.TrimSpace(c.Param("subscriptionID"))
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscriptionID is required"})
		return
	}

	sub, err := internal.GetSubscriptionByID(subscriptionID)
	if err != nil {
		log.Printf("poll subscription %s: %v", subscriptionID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}
