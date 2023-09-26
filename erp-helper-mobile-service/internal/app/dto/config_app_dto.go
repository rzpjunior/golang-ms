package dto

type ConfigAppResponse struct {
	Id          int64  `json:"id"`
	Application int8   `json:"application"`
	Field       string `json:"field"`
	Attribute   string `json:"attribute"`
	Value       string `json:"value"`
}

// Get ConfigApp
type GetConfigAppRequest struct {
	Limit       int
	Offset      int
	Id          int
	Application int
	Field       string
	Attribute   string
	Value       string
}
