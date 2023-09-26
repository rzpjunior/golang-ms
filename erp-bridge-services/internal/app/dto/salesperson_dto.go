package dto

import "time"

type SalespersonResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	FirstName     string    `json:"firstname"`
	MiddleName    string    `json:"namemiddle"`
	LastName      string    `json:"lastname"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
