package models

type BillingCycle string

const (
	BillingCycleMonthly  BillingCycle = "monthly"
	BillingCycleYearly   BillingCycle = "yearly"
	BillingCycleLifeTime BillingCycle = "life_time"
)

type Plan struct {
	ID           string       `json:"id"`
	CreatedAt    string       `json:"created_at"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	BillingCycle BillingCycle `json:"billing_cycle"`
}

type VideoPlan struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	VideoID   string `json:"video_id"`
	PlanID    string `json:"plan_id"`
}
