package dto

type CustomerTypeResponse struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	CustomerGroup string `json:"customer_group,omitempty"`
	Status        int8   `json:"status,omitempty"`
	ConvertStatus string `json:"convert_status,omitempty"`
}

type CustomerTypeGetListRequest struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Search string `json:"search"`
	Status int8   `json:"status"`
}
