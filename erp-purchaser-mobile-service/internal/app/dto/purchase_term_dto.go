package dto

type PurchaseTermResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	DaysValue   int64  `json:"days_value"`
	Note        string `json:"note,omitempty"`
	Status      int8   `json:"status"`
}
