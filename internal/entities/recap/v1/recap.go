package entities_recap_v1

import (
	entities_payments_v1 "github.com/subeecore/subee-core-svc/internal/entities/payments/v1"
)

type MonthlyRecap struct {
	Price    float64                         `json:"price"`
	Payments []*entities_payments_v1.Payment `json:"payments"`
}
