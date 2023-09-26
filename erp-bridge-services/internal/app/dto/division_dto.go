package dto

type DivisionResponse struct {
	ID            int64  `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Note          string `json:"note"`
	Status        int8   `json:"status"`
	StatusConvert string `json:"status_convert"`
}
