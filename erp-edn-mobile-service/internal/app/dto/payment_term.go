package dto

type PaymentTermListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type PaymentTermGP struct {
	ID                string `json:"id"`
	PaymentTermCode   string `json:"code"`
	CalculateFromDays string `json:"days_value"`
}
