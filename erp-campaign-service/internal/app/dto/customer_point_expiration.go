package dto

import "time"

type CustomerPointExpirationResponse struct {
	ID                 int64     `json:"id"`
	CustomerID         int64     `json:"customer_id"`
	CurrentPeriodPoint float64   `json:"current_period_point"`
	NextPeriodPoint    float64   `json:"next_period_point"`
	CurrentPeriodDate  time.Time `json:"current_period_date"`
	NextPeriodDate     time.Time `json:"next_period_date"`
	LastUpdatedAt      time.Time `json:"last_updated_at"`
}
