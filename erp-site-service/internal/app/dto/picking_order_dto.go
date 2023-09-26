package dto

type GetListProduct struct {
	DocNumber string              `json:"doc_number"`
	Status    int8                `json:"status"`
	Products  []*AggregateProduct `json:"products"`
}

type AggregateProduct struct {
	ItemNumber      string  `json:"item_number"`
	ItemName        string  `json:"item_name"`
	Picture         string  `json:"picture"`
	UomDescription  string  `json:"uom_description"`
	TotalOrderQty   float64 `json:"total_order_qty"`
	TotalPickedQty  float64 `json:"total_picked_qty"`
	TotalSalesOrder int32   `json:"total_sales_order"`
	Status          int8    `json:"status"` // 1 = none , 2 = picked , 3 = rejected
}
