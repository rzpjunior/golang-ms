package dto

type AdmDivisionResponse struct {
	Code        string `json:"code"`
	Region      string `json:"region"`
	State       string `json:"state"`
	City        string `json:"city"`
	District    string `json:"district"`
	SubDistrict string `json:"sub_district"`
}
