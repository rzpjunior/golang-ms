package dto

type DeliveryFeeResponse struct {
	ID            int64   `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name,omitempty"`
	Note          string  `json:"note,omitempty"`
	Status        int32   `json:"status,omitempty"`
	MinimumOrder  float64 `json:"minimum_order,omitempty"`
	DeliveryFee   float64 `json:"delivery_fee,omitempty"`
	RegionId      int64   `json:"region_id,omitempty"`
	CutomerTypeId string  `json:"customer_type_id,omitempty"`
}
