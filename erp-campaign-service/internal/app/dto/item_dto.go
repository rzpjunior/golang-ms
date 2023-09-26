package dto

type ItemResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Status      int8   `json:"status,omitempty"`
}
