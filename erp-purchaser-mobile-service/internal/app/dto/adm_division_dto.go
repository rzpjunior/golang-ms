package dto

type AdmDivisionResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Region      string `json:"region"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	SubDistrict string `json:"sub_district"`
}
