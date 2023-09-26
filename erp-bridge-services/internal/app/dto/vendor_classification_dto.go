package dto

type VendorClassificationResponse struct {
	ID            int64  `json:"id"`
	CommodityCode string `json:"commodity_code"`
	CommodityName string `json:"commodity_name"`
	BadgeCode     string `json:"badge_code"`
	BadgeName     string `json:"badge_name"`
	TypeCode      string `json:"type_code"`
	TypeName      string `json:"type_name"`
}
