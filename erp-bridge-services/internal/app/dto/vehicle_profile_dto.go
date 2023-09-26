package dto

type VehicleProfileResponse struct {
	ID                  int64   `json:"id"`
	Code                string  `json:"code"`
	Name                string  `json:"name"`
	MaxKoli             float64 `json:"max_koli"`
	MaxWeight           float64 `json:"max_weight"`
	MaxFragile          float64 `json:"max_fragile"`
	SpeedFactor         float64 `json:"speed_factor"`
	RoutingProfile      int8    `json:"routing_profile"`
	Status              int8    `json:"status"`
	StatusConvert       string  `json:"status_convert"`
	Skills              string  `json:"skills"`
	InitialCost         float64 `json:"initial_cost"`
	SubsequentCost      float64 `json:"subsequent_cost"`
	MaxAvailableVehicle int64   `json:"max_available_vehicle"`

	CourierVendorID int64 `json:"courier_vendor_id"`
}
