package dto

import "time"

type CustomerTypeResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	GroupType     string    `json:"group_type"`
	Abbreviation  string    `json:"abbreviation"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
