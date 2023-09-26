package dto

import "time"

type AdmDivisionResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	ProvinceID    int64     `json:"province_id"`
	CityID        int64     `json:"city_id"`
	DistrictID    int64     `json:"district_id"`
	SubDistrictID int64     `json:"sub_district_id"`
	RegionID      int64     `json:"region_id"`
	PostalCode    string    `json:"postal_code"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	Region        string    `json:"region"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
