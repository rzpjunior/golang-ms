package dto

import "time"

type ArchetypeResponse struct {
	ID             int64                 `json:"id"`
	Code           string                `json:"code"`
	CustomerTypeID int64                 `json:"customer_type_id"`
	Description    string                `json:"description"`
	Status         int8                  `json:"status"`
	StatusConvert  string                `json:"status_convert"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	CustomerType   *CustomerTypeResponse `json:"customer_type"`
}
