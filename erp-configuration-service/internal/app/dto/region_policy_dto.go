package dto

type RegionPolicyResponse struct {
	ID                 int64               `json:"id,omitempty"`
	OrderTimeLimit     string              `json:"order_time_limit"`
	MaxDayDeliveryDate int                 `json:"max_day_delivery_date"`
	WeeklyDayOff       int                 `json:"weekly_day_off"`
	CSPhoneNumber      string              `json:"cs_phone_number"`
	Region             *RegionResponse     `json:"region,omitempty"`
	DefaultPriceLevel  *PriceLevelResponse `json:"default_price_level"`
}

type RegionPolicyRequestUpdate struct {
	OrderTimeLimit     string `json:"order_time_limit" valid:"required"`
	MaxDayDeliveryDate int    `json:"max_day_delivery_date" valid:"required"`
	WeeklyDayOff       int    `json:"weekly_day_off" valid:"required"`
	DefaultPriceLevel  string `json:"default_price_level" valid:"required"`
}

type PriceLevelResponse struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	CustomerTypeID string `json:"customer_type_id"`
	RegionID       string `json:"region_id"`
}
