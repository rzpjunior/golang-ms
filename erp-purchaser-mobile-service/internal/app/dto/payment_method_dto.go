package dto

type PaymentMethodResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Status      int8   `json:"status"`
	Publish     int8   `json:"publish"`
	Maintenance int8   `json:"maintenance"`
}

type PaymentMethodListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}
