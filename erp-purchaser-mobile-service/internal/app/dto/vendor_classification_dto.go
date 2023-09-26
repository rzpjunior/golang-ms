package dto

type VendorClassificationResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type VendorClassificationListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}
