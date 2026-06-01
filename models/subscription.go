package models

type SubscriptionStatus string

const (
	SubscriptionStatusPending  SubscriptionStatus = "pending"
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusFailed   SubscriptionStatus = "failed"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
)

type Subscription struct {
	ID        int8   `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UserID    string `json:"user_id"`
	PlanID    string `json:"plan_id"`
	StartedAt string `json:"started_at,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
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
