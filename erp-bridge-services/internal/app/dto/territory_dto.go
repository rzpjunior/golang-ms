package dto

import "time"

type TerritoryResponse struct {
	ID             int64     `json:"id"`
	Code           string    `json:"code"`
	Description    string    `json:"description"`
	RegionID       int64     `json:"region_id"`
	SalespersonID  int64     `json:"salesperson_id"`
	CustomerTypeID int64     `json:"customer_type_id"`
	SubDistrictID  int64     `json:"sub_district_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
