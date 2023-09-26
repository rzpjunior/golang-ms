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

type AdmDivisionListRequest struct {
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
	Status      int32  `json:"status"`
	Search      string `json:"search"`
	OrderBy     string `json:"order_by"`
	Region      string `json:"region"`
	Code        string `json:"code"`
	State       string `json:"state"`
	Type        string `json:"type"`
	City        string `json:"city"`
	District    string `json:"district"`
	Subdistrict string `json:"subdistrict"`
}

type AdmDivisionDetailRequest struct {
	Id int32 `json:"id"`
}

type AdmDivisionGP struct {
	Code        string `json:"code,omitempty"`
	Region      string `json:"region,omitempty"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	District    string `json:"district,omitempty"`
	SubDistrict string `json:"sub_district,omitempty"`
}

type AdmDivisionCoverageListRequest struct {
	Limit                 int32  `json:"limit"`
	Offset                int32  `json:"offset"`
	Status                int32  `json:"status"`
	Code                  string `json:"code"`
	OrderBy               string `json:"order_by"`
	GnlProvince           string `json:"gnl_province"`
	GnlAdministrativeCode string `json:"gnl_administrative_code"`
	GnlCity               string `json:"gnl_city"`
	GnlDistrict           string `json:"gnl_district"`
	GnlSubdistrict        string `json:"gnl_subdistrict"`
	Locncode              string `json:"locncode"`
}

type AdmDivisionCoverageGP struct {
	GnlAdministrativeCode string `json:"gnl_administrative_code"`
	GnlRegion             string `json:"gnl_region"`
	GnlProvince           string `json:"gnl_province"`
	GnlCity               string `json:"gnl_city"`
	GnlDistrict           string `json:"gnl_district"`
	GnlSubdistrict        string `json:"gnl_subdistrict"`
	Locncode              string `json:"locncode"`
}

type AdmDivisionGPResponse struct {
	Code          string `json:"code"`
	Region        string `json:"region"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	SubDistrict   string `json:"sub_district"`
	ZipCode       string `json:"zip_code"`
	ConcatAddress string `json:"concat_address,omitempty"`
}

type GetAdmDivisionRequest struct {
	Region      string `json:"region"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	SubDistrict string `json:"sub_district"`
	TypeAdm     string `json:"type_adm"`
}
