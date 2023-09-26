package dto

type DistrictResponse struct {
	ID            int64  `json:"-"`
	Code          string `json:"code"`
	Value         string `json:"value"`
	Name          string `json:"name"`
	Note          string `json:"note"`
	Status        int8   `json:"status"`
	StatusConvert string `json:"status_convert"`
}

// TODO: Connect to GP
