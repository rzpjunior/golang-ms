package dto

type SalesPriceLevelResponse struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	CustomerTypeID string `json:"customer_type_id"`
	RegionID       string `json:"region_id"`
}

type GetSalesPriceLevelRequest struct {
	Limit          int64
	Offset         int64
	RegionID       string
	CustomerTypeID string
}
