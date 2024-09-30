package entities_subscriptions_v1

import "time"

type Subscription struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Platform   string     `json:"platform"`
	Reccurence int        `json:"reccurence"`
	Price      float64    `json:"price"`
	StartedAt  time.Time  `json:"started_at"`
	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	UserID     string  `json:"user_id"`
	Platform   string  `json:"platform"`
	Reccurence int     `json:"reccurence"`
	Price      float64 `json:"price"`
	StartedAt  string  `json:"started_at"`
}
