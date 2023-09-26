package dto

type ApplicationConfigResponse struct {
	ID          int64  `json:"id,omitempty"`
	Application int8   `json:"application"`
	Field       string `json:"field"`
	Attribute   string `json:"attribute"`
	Value       string `json:"value"`
}

type ApplicationConfigRequestCreate struct {
	Application int8   `json:"application"`
	Field       string `json:"field"`
	Attribute   string `json:"attribute"`
	Value       string `json:"value"`
}

type ApplicationConfigRequestUpdate struct {
	Application int8   `json:"application"`
	Field       string `json:"field"`
	Value       string `json:"value" valid:"required"`
}

type ApplicationConfigRequestGet struct {
	Offset      int32
	Limit       int32
	Status      int8
	Search      string
	OrderBy     string
	Application int8
	Attribute   string
	Value       string
}
