package dto

type CustomerClassRequest struct {
	Limit          int64  `json:"limit"`
	Offset         int64  `json:"offset"`
	Search         string `json:"search"`
	OrderBy        string `json:"order_by"`
}

type CustomerClassResponse struct {
	ID                  string `json:"id"`
	Description         string `json:"description,omitempty"`
	CreditLimitType     int32 `json:"credit_limit_type"`
	CreditLimitTypeDesc string `json:"credit_limit_type_desc"`
	CreditLimitAmount   float64 `json:"credit_limit_amount"`
}
