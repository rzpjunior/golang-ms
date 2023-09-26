package dto

type AdmDivisionResponse struct {
	Region      string `json:"region"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	SubDistrict string `json:"sub_district"`
	PostalCode  string `json:"postal_code"`
}

type AdmDivisionGetRequest struct {
	Region            string `json:"region"`
	RegionSearch      string `json:"region_search"`
	Province          string `json:"province"`
	ProvinceSearch    string `json:"province_search"`
	City              string `json:"city"`
	CitySearch        string `json:"city_search"`
	District          string `json:"district"`
	DistrictSearch    string `json:"district_search"`
	SubDistrict       string `json:"sub_district"`
	SubDistrictSearch string `json:"sub_district_search"`
	Limit             int64  `json:"limit"`
	Offset            int64  `json:"offset"`
}
