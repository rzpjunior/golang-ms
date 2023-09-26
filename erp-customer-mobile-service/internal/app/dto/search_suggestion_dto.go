package dto

type SearchSuggestionResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type SearchSuggestionRequest struct {
	Platform string     `json:"platform" valid:"required"`
	Data     dataSearch `json:"data" valid:"required"`
}

type dataSearch struct {
	Search string `json:"search"`
}