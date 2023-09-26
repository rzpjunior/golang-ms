package dto

type ShippingMethodResponse struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	Type            int8   `json:"type"`
	TypeDescription string `json:"type_description"`
}

type GetShippingMethodRequest struct {
	Limit  int64
	Offset int64
	Type   string
}
