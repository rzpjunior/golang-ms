package dto

import "time"

type BankResponse struct {
	PaymentMethodName string `json:"payment_method_name"`
	ListBank          []Bank `json:"list_bank"`
}

type Bank struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	ImageUrl        string `json:"image_url"`
	PaymentGuideUrl string `json:"payment_guide_url"`
}

type InvoiceXenditModify struct {
	ServerTime      time.Time `json:"server_time"`
	DeadlinePayment time.Time `json:"deadline_payment"`
	VaNumber        string    `json:"va_number"`
	PaymentNominal  float64   `json:"payment_nominal"`
	ImageUrl        string    `json:"image_url"`
	PaymentGuideUrl string    `json:"payment_guide_url"`
	Name            string    `json:"name"`
}

type PaymentOption struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	ImageURL        string `json:"image_url"`
	PaymentGuideURL string `json:"payment_guide_url"`
}

type PaymentMethod struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Value          string           `json:"value"`
	Note           string           `json:"note"`
	PaymentOptions []*PaymentOption `json:"payment_options,omitempty"`
}

type PaymentMethodRequestGet struct {
	Session *SessionDataCustomer
}

type PaymentRequest struct {
	Platform     string `json:"platform" valid:"required"`
	SalesOrderID int64  `json:"sales_order_id" valid:"required"`
	Session      *SessionDataCustomer
}
