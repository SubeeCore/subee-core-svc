package entities_payments_v1

type Payment struct {
	SubscriptionID string  `json:"subscription_id"`
	Platform       string  `json:"platform"`
	Price          float64 `json:"price"`
}
