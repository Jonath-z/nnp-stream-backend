package models

type SubscriptionStatus string

const (
	SubscriptionStatusPending  SubscriptionStatus = "pending"
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusFailed   SubscriptionStatus = "failed"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
)

type Subscription struct {
	ID        string             `json:"id,omitempty"`
	CreatedAt string             `json:"created_at,omitempty"`
	UserID    string             `json:"user_id"`
	PlanID    string             `json:"plan_id"`
	Status    SubscriptionStatus `json:"status,omitempty"`
	StartedAt string             `json:"start_at,omitempty"`
	ExpiresAt string             `json:"end_at,omitempty"`
}

type PaymentTransaction struct {
	ID             string  `json:"id,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UserID         string  `json:"user_id"`
	PlanID         string  `json:"plan_id"`
	SubscriptionID string  `json:"subscription_id,omitempty"`
	Provider       string  `json:"provider"`
	ProviderTxID   string  `json:"provider_tx_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	PhoneNumber    string  `json:"phone_number"`
	Status         string  `json:"status"`
	FailureReason  string  `json:"failure_reason,omitempty"`
	TxHash         string  `json:"tx_hash,omitempty"`
	CompletedAt    string  `json:"completed_at,omitempty"`
}

type SubscribeRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	PlanID      string  `json:"plan_id" binding:"required"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}

type VideoAccessReason string

const (
	VideoAccessReasonActive         VideoAccessReason = "active"
	VideoAccessReasonPending        VideoAccessReason = "pending"
	VideoAccessReasonExpired        VideoAccessReason = "expired"
	VideoAccessReasonNoSubscription VideoAccessReason = "no_subscription"
	VideoAccessReasonFreeVideo      VideoAccessReason = "free_video"
)

type CheckVideoAccessRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	VideoID string `json:"video_id" binding:"required"`
}

type CheckVideoAccessResponse struct {
	HasAccess    bool              `json:"has_access"`
	Reason       VideoAccessReason `json:"reason"`
	Subscription *Subscription     `json:"subscription,omitempty"`
}
