package dto

import "time"

type RegionResponse struct {
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Code          string    `json:"code,omitempty"`
	Description   string    `json:"description,omitempty"`
	Status        int8      `json:"status,omitempty"`
	StatusConvert string    `json:"status_convert,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type RegionListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type RegionDetailRequest struct {
	Id int32 `json:"id"`
}
