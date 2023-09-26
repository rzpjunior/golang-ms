package dto

type PaymentMethodListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type PaymentMethodGP struct {
	ID                string `json:"id"`
	PaymentMethodCode string `json:"payment_method_code"`
	PaymentMethodDesc string `json:"payment_method_desc"`
}
