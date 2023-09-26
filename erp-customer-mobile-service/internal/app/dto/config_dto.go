package dto

type ApplicationConfigResponse struct {
	ID          string `json:"id,omitempty"`
	Application string `json:"application"`
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
	Attribute   string `json:"attribute"`
	Value       string `json:"value" valid:"required"`
}

type GlossaryResponse struct {
	ID        string `json:"id,omitempty"`
	Table     string `json:"table,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	ValueInt  string `json:"value_int"`
	ValueName string `json:"value_name"`
	Note      string `json:"note"`
}

type RequestGetDeliveryFee struct {
	Platform string             `json:"platform" valid:"required"`
	Data     dataGetDeliveryFee `json:"data" valid:"required"`
}

type dataGetDeliveryFee struct {
	RegionID       string `json:"region_id" valid:"required"`
	CustomerTypeID string `json:"customer_type_id" valid:"required"`
	// DataResponse   responseGetData

	// Area         *model.Area
	// BusinessType *model.BusinessType
}

type ResponseGetDeliveryFee struct {
	ID          string `json:"id"`
	MinOrder    string `json:"min_order"`
	DeliveryFee string `json:"delivery_fee"`
}
