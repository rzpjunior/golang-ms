package dto

import "time"

type OrderTypeResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
