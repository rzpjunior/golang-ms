package dto

type PaymentTermResponse struct {
	ID                       string `json:"id"`
	Description              string `json:"description"`
	DueType                  string `json:"due_type"`
	PaymentUseFor            int    `json:"payment_usefor"`
	PaymentUseForDescription string `json:"payment_usefor_description"`
}

type GetPaymentTermRequest struct {
	Limit         int
	Offset        int
	PaymentUseFor string
}
