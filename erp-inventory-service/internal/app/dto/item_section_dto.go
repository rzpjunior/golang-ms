package dto

type ItemSectionResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
