package entities_recap_v1

import (
	entities_categories_v1 "github.com/subeecore/subee-core-svc/internal/entities/categories/v1"
	entities_payments_v1 "github.com/subeecore/subee-core-svc/internal/entities/payments/v1"
)

type MonthlyRecap struct {
	Price      float64                                 `json:"price"`
	Payments   []*entities_payments_v1.Payment         `json:"payments"`
	Categories *entities_categories_v1.CategoriesRecap `json:"categories"`
}
