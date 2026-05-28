package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nnp-stream-backend/models"
)

const (
	TABLE_PLANS         = "plans"
	TABLE_SUBSCRIPTIONS = "subscriptions"
	TABLE_TRANSACTIONS  = "payment_transactions"
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
