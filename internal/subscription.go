package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nnp-stream-backend/models"
)

const (
	TABLE_PLANS         = "plans"
	TABLE_SUBSCRIPTIONS = "user_subscriptions"
	TABLE_TRANSACTIONS  = "payment_transactions"
	TABLE_VIDEO_PLANS   = "video_plans"
)

func GetPlan(planID string) (models.Plan, error) {
	var plan models.Plan
	client, err := SupabaseClient()
	if err != nil {
		return plan, fmt.Errorf("init supabase: %w", err)
	}

	raw, _, err := client.From(TABLE_PLANS).Select("*", "exact", false).Eq("id", planID).Single().Execute()
	if err != nil {
		return plan, fmt.Errorf("fetch plan: %w", err)
	}
	if err := json.Unmarshal(raw, &plan); err != nil {
		return plan, fmt.Errorf("decode plan: %w", err)
	}
	return plan, nil
}

func CreateSubscription(sub models.Subscription) (models.Subscription, error) {
	var created models.Subscription
	client, err := SupabaseClient()
	if err != nil {
		return created, fmt.Errorf("init supabase: %w", err)
	}

	raw, _, err := client.From(TABLE_SUBSCRIPTIONS).Insert(&sub, false, "", "representation", "").Execute()
	if err != nil {
		return created, fmt.Errorf("insert subscription: %w", err)
	}

	var rows []models.Subscription
	if err := json.Unmarshal(raw, &rows); err != nil {
		return created, fmt.Errorf("decode subscription: %w", err)
	}
	if len(rows) == 0 {
		return created, fmt.Errorf("empty insert response")
	}
	return rows[0], nil
}

func UpdateSubscriptionStatus(subscriptionID string, status models.SubscriptionStatus, expiresAt string) error {
	client, err := SupabaseClient()
	if err != nil {
		return fmt.Errorf("init supabase: %w", err)
	}

	patch := map[string]interface{}{"status": status}
	if status == models.SubscriptionStatusActive {
		patch["started_at"] = time.Now().UTC().Format(time.RFC3339)
		if expiresAt != "" {
			patch["expires_at"] = expiresAt
		}
	}

	_, _, err = client.From(TABLE_SUBSCRIPTIONS).Update(patch, "", "").Eq("id", subscriptionID).Execute()
	if err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}
	return nil
}

func CreatePaymentTransaction(tx models.PaymentTransaction) (models.PaymentTransaction, error) {
	var created models.PaymentTransaction
	client, err := SupabaseClient()
	if err != nil {
		return created, fmt.Errorf("init supabase: %w", err)
	}

	raw, _, err := client.From(TABLE_TRANSACTIONS).Insert(&tx, false, "", "representation", "").Execute()
	if err != nil {
		return created, fmt.Errorf("insert transaction: %w", err)
	}

	var rows []models.PaymentTransaction
	if err := json.Unmarshal(raw, &rows); err != nil {
		return created, fmt.Errorf("decode transaction: %w", err)
	}
	if len(rows) == 0 {
		return created, fmt.Errorf("empty insert response")
	}
	return rows[0], nil
}

func GetPaymentTransactionByProviderID(provider, providerTxID string) (models.PaymentTransaction, error) {
	var tx models.PaymentTransaction
	client, err := SupabaseClient()
	if err != nil {
		return tx, fmt.Errorf("init supabase: %w", err)
	}

	raw, _, err := client.From(TABLE_TRANSACTIONS).
		Select("*", "exact", false).
		Eq("provider", provider).
		Eq("provider_tx_id", providerTxID).
		Single().
		Execute()
	if err != nil {
		return tx, fmt.Errorf("fetch transaction: %w", err)
	}
	if err := json.Unmarshal(raw, &tx); err != nil {
		return tx, fmt.Errorf("decode transaction: %w", err)
	}
	return tx, nil
}

func UpdatePaymentTransaction(transactionID string, patch map[string]interface{}) error {
	client, err := SupabaseClient()
	if err != nil {
		return fmt.Errorf("init supabase: %w", err)
	}
	_, _, err = client.From(TABLE_TRANSACTIONS).Update(patch, "", "").Eq("id", transactionID).Execute()
	if err != nil {
		return fmt.Errorf("update transaction: %w", err)
	}
	return nil
}

func GetSubscriptionByID(id string) (models.Subscription, error) {
	var sub models.Subscription
	client, err := SupabaseClient()
	if err != nil {
		return sub, fmt.Errorf("init supabase: %w", err)
	}
	raw, _, err := client.From(TABLE_SUBSCRIPTIONS).Select("*", "exact", false).Eq("id", id).Single().Execute()
	if err != nil {
		return sub, fmt.Errorf("fetch subscription: %w", err)
	}
	if err := json.Unmarshal(raw, &sub); err != nil {
		return sub, fmt.Errorf("decode subscription: %w", err)
	}
	return sub, nil
}

// GetPlanIDsForVideo returns the plan IDs that gate access to the given video.
// An empty slice means the video is not gated by any paid plan.
func GetPlanIDsForVideo(videoID string) ([]string, error) {
	client, err := SupabaseClient()
	if err != nil {
		return nil, fmt.Errorf("init supabase: %w", err)
	}
	raw, _, err := client.From(TABLE_VIDEO_PLANS).Select("plan_id", "exact", false).Eq("video_id", videoID).Execute()
	if err != nil {
		return nil, fmt.Errorf("fetch video plans: %w", err)
	}
	var rows []struct {
		PlanID string `json:"plan_id"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("decode video plans: %w", err)
	}
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.PlanID)
	}
	return ids, nil
}

// FindUserSubscriptionForPlans returns the user's most relevant subscription
// among the given plan IDs. It prefers an active, unexpired subscription;
// otherwise it returns the most recent pending one (so the caller can report
// "payment still processing"). Returns (nil, nil) when no subscription matches.
func FindUserSubscriptionForPlans(userID string, planIDs []string) (*models.Subscription, error) {
	if len(planIDs) == 0 {
		return nil, nil
	}
	client, err := SupabaseClient()
	if err != nil {
		return nil, fmt.Errorf("init supabase: %w", err)
	}
	raw, _, err := client.From(TABLE_SUBSCRIPTIONS).
		Select("*", "exact", false).
		Eq("user_id", userID).
		In("plan_id", planIDs).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("fetch subscriptions: %w", err)
	}
	var subs []models.Subscription
	if err := json.Unmarshal(raw, &subs); err != nil {
		return nil, fmt.Errorf("decode subscriptions: %w", err)
	}

	now := time.Now().UTC()
	var pending *models.Subscription
	for i := range subs {
		s := subs[i]
		if s.Status == models.SubscriptionStatusActive && !isExpired(s.ExpiresAt, now) {
			return &s, nil
		}
		if s.Status == models.SubscriptionStatusPending && pending == nil {
			pending = &s
		}
	}
	return pending, nil
}

func isExpired(expiresAt string, now time.Time) bool {
	if expiresAt == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return false
	}
	return !t.After(now)
}

// ComputeSubscriptionExpiry returns the expiry timestamp for a given plan's billing cycle.
// life_time plans return an empty string (no expiry).
func ComputeSubscriptionExpiry(cycle models.BillingCycle, from time.Time) string {
	switch cycle {
	case models.BillingCycleMonthly:
		return from.AddDate(0, 1, 0).UTC().Format(time.RFC3339)
	case models.BillingCycleYearly:
		return from.AddDate(1, 0, 0).UTC().Format(time.RFC3339)
	}
	return ""
}
