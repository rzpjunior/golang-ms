package dto

type ConfigAppResponse struct {
	ID          int64  `json:"id,omitempty"`
	Application int8   `json:"application"`
	Field       string `json:"field"`
	Attribute   string `json:"attribute"`
	Value       string `json:"value"`
}
