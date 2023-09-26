package dto

type CustomerClassResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type CustomerClassGetListRequest struct {
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
	Search string `json:"search"`
}
