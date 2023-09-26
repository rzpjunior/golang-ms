package dto

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

type PaymentTermResponse struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"code,omitempty"`
}
