package dto

type UomResponse struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name,omitempty"`
	DecimalEnabled int    `json:"decimal_enabled"`
}
