package dto

type VendorOrganizationResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type VendorOrganizationListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}
