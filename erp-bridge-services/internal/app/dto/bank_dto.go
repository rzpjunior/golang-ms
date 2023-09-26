package dto

import "time"

type BankResponse struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	Description     string    `json:"description"`
	Value           string    `json:"value"`
	ImageUrl        string    `json:"image_url"`
	PaymentGuideUrl string    `json:"payment_guide_url"`
	PublishIVA      int8      `json:"publish_iva"`
	PublishFVA      int8      `json:"publish_fva"`
	Status          int8      `json:"status"`
	StatusConvert   string    `json:"status_convert"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
