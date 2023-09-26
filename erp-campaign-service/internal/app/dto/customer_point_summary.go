package dto

type CustomerPointSummaryResponse struct {
	ID            int64   `json:"id"`
	CustomerID    int64   `json:"customer_id"`
	EarnedPoint   float64 `json:"earned_point"`
	RedeemedPoint float64 `json:"redeemed_point"`
	SummaryDate   string  `json:"summary_date"`
}

type CustomerPointSummaryRequestCreate struct {
	CustomerID    int64   `json:"customer_id"`
	EarnedPoint   float64 `json:"earned_point"`
	RedeemedPoint float64 `json:"redeemed_point"`
	SummaryDate   string  `json:"summary_date"`
}

type CustomerPointSummaryRequestUpdate struct {
	ID            int64    `json:"id"`
	CustomerID    int64    `json:"customer_id"`
	EarnedPoint   float64  `json:"earned_point"`
	RedeemedPoint float64  `json:"redeemed_point"`
	FieldUpdate   []string `json:"field_update"`
}

type CustomerPointSummaryRequestGetDetail struct {
	ID          int64  `json:"id"`
	CustomerID  int64  `json:"customer_id"`
	SummaryDate string `json:"summary_date"`
}
