package dto

type RegionPolicy struct {
	RegPol RegionPolicyResponse `json:"RegionPolicy"`
}
type RegionPolicyResponse struct {
	ID                 string          `json:"id,omitempty"`
	OrderTimeLimit     string          `json:"order_time_limit"`
	MaxDayDeliveryDate string          `json:"max_day_delivery_date"`
	WeeklyDayOff       string          `json:"weekly_day_off"`
	CSPhoneNumber      string          `json:"cs_phone_number"`
	Region             *RegionResponse `json:"region,omitempty"`
}

type RegionPolicyRequestUpdate struct {
	OrderTimeLimit     string `json:"order_time_limit" valid:"required"`
	MaxDayDeliveryDate int    `json:"max_day_delivery_date" valid:"required"`
	WeeklyDayOff       int    `json:"weekly_day_off" valid:"required"`
}

type RegionPolicyMobileRequest struct {
	Platform string        `json:"platform" valid:"required"`
	Data     DataGetMobile `json:"data" valid:"required"`
}

type DataGetMobile struct {
	AdmDivisionID    string `json:"adm_division_id" valid:"required"`
	PaymentGroupCode string `json:"payment_group_code"`
}
