package dto

import "time"

type UomResponse struct {
	ID             string    `json:"id"`
	Code           string    `json:"code"`
	Description    string    `json:"description"`
	Status         int8      `json:"status"`
	StatusConvert  string    `json:"status_convert"`
	DecimalEnabled int8      `json:"decimal_enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
