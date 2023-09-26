package dto

import "time"

type AdmDivisionResponse struct {
	ID              string    `json:"id"`
	Code            string    `json:"code"`
	RegionID        string    `json:"region_id"`
	RegionName      string    `json:"region_name"`
	Province        string    `json:"province"`
	City            string    `json:"city"`
	District        string    `json:"district"`
	SubDistrictID   string    `json:"sub_district_id"`
	SubDistrictName string    `json:"sub_district_name"`
	DistrictID      string    `json:"district_id"`
	DistrictName    string    `json:"district_name"`
	CityID          string    `json:"city_id"`
	CityName        string    `json:"city_name"`
	ProvinceID      string    `json:"province_id"`
	ProvinceName    string    `json:"province_name"`
	ContryID        string    `json:"country_id"`
	CountryName     string    `json:"country_name"`
	PostalCode      string    `json:"postal_code"`
	Status          string    `json:"status"`
	StatusConvert   string    `json:"status_convert"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AdmDivisionGPResponse struct {
	Code          string `json:"code"`
	Region        string `json:"region"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	SubDistrict   string `json:"sub_district"`
	ConcatAddress string `json:"concat_address,omitempty"`
}
type GetAdmDivisionRequest struct {
	Platform string         `json:"platform" valid:"required"`
	Data     GetAdmDivision `json:"data"`
}

type GetAdmDivision struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
}

type SearchAdmDivisionRequest struct {
	Platform string            `json:"platform" valid:"required"`
	Data     SearchAdmDivision `json:"data"`
}

type SearchAdmDivision struct {
	SubDistrict string `json:"sub_district"`
}
