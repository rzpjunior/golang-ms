package dto

import "time"

type PaymentTermResponse struct {
	Id          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      int32     `json:"status,omitempty"`
	DaysValue   int64     `json:"days_value"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type PaymentTermListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}
