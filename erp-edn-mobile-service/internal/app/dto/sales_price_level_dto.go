package dto

type GetSalesPriceLevelGPListRequest struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	GnlRegion     string `json:"gnl_region"`
	GnlCustTypeId string `json:"gnl_cust_type_id"`
	Prclevel      string `json:"prclevel"`
}

type SalesPriceLevelGP struct {
	GnlRegion     string `json:"gnl_region"`
	GnlCustTypeId string `json:"gnl_cust_type_id"`
	Prclevel      string `json:"prclevel"`
}

type GetSalesPriceLevelListRequest struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	RegionID   string `json:"region_id"`
	CustTypeID string `json:"customer_type_id"`
	PriceLevel string `json:"price_level"`
}
type SalesPriceLevel struct {
	RegionID   string `json:"region_id"`
	CustTypeID string `json:"customer_type_id"`
	PriceLevel string `json:"price_level"`
}
