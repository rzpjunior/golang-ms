package dto

import "time"

type SubDistrictResponse struct {
	ID            int64                `json:"id"`
	Code          string               `json:"code"`
	Description   string               `json:"description"`
	District      *AdmDivisionResponse `json:"district"`
	Status        int8                 `json:"status"`
	StatusConvert string               `json:"status_convert"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}
