package entities_payments_v1

import "time"

type Payment struct {
	SubscriptionID string    `json:"subscription_id"`
	Platform       string    `json:"platform"`
	Category       string    `json:"category"`
	Price          float64   `json:"price"`
	Reccurence     int       `json:"reccurence"`
	StartedAt      time.Time `json:"started_at"`
}

type Payment_Light struct {
	Price     float64   `json:"price"`
	StartedAt time.Time `json:"started_at"`
}
